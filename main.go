// Simple REST API with a PostgreSQL backend
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type InboundPayload struct {
	Time	string 	`json:"time"`
	Key	string 	`json:"key"`
	Field1	string 	`json:"field1"`
	Field2	string 	`json:"field2"`
}

type Config struct {
	Host	string  `json:"host"`
	Port	int     `json:"port"`
	User	string  `json:"user"`
	Passwd  string  `json:"passwd"`
	DBName  string  `json:"dbname"`
	DBPort  int     `json:"dbport"`
	Sslmode	string  `json:"sslmode"`
}

var (
	conf	= flag.String("conf", "rest_api.conf", "App configuration data")
	logFile = flag.String("log", "rest_api.log", "Location for app log file")
	config	  Config
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func root(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp := make(map[string]string)
	resp["status"] = "ok"
	resp["message"] = "success"
	log.Println("root: ", resp["message"])
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func echo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp := make(map[string]string)
	resp["status"] = "ok"
	resp["message"] = "echo"
	log.Println(resp["message"])
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func endpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp := make(map[string]string)
	var data InboundPayload

	if r.Method == "POST" {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["status"] = "not_ok"
			resp["message"] = err.Error()
			log.Println("io.ReadAll() error:", resp["message"])
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w, "{\"status\":\"not_ok\",\"message\":\"json.Marshal() error\"}")
				log.Println("json.Marshal() error", err)
				return
			}
			w.Write(jsonResp)
			return
		}

		err = json.Unmarshal(reqBody, &data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["status"] = "not_ok"
			resp["message"] = err.Error()
			log.Println("json.Unmarshal() errpr:", resp["message"])
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w, "{\"status\":\"not_ok\",\"message\":\"json.Marshal() error\"}")
				log.Println("json.Marshal() error:", err)
				return
			}
			w.Write(jsonResp)
			return
		}

		connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
					config.Host, config.DBPort, config.User, config.Passwd, config.DBName, config.Sslmode)

		conn, err := sql.Open("postgres", connString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["status"] = "not_ok"
			resp["message"] = err.Error()
			log.Println("sql.Open() error:", resp["message"])
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w, "{\"status\":\"not_ok\",\"message\":\"json.Marshal() error\"}")
				log.Println("json.Marshal() error", err)
				return
			}
			w.Write(jsonResp)
			return
		}

		defer conn.Close()
		statement := fmt.Sprintf("insert into test (time,key,field1,field2) values ('%s','%s','%s','%s');",
					data.Time, data.Key, data.Field1, data.Field2)

		_, err = conn.Exec(statement)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["status"] = "not_ok"
			resp["message"] = err.Error()
			log.Println("conn.Exec() error:", resp["message"])
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w, "{\"status\":\"not_ok\",\"message\":\"json.Marshal() error\"}")
				log.Println("conn.Exec(): json.Marshal() error: ", err)
				return
			}
			w.Write(jsonResp)
			return
		}
		resp["status"] = "ok"
		resp["message"] = "created"
		jsonResp, err := json.Marshal(resp)
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResp)
	} else if r.Method == "GET" {
		w.WriteHeader(http.StatusNotFound)
		resp["status"] = "ok"
		resp["message"] = "method not implemented"
		log.Println(resp["message"])
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "{\"status\":\"not_ok\",\"message\":\"json.Marshal() fail\"}")
			log.Println("GET json.Marshal() fail: ", err)
			return
		}
		w.Write(jsonResp)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp["status"] = "ok"
		resp["message"] = "method not allowed"
		log.Println(r.Method, "not allowed")
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "{\"status\":\"not_ok\",\"message\":\"json.Marshal() fail\"}")
			log.Println("Method Not Allowed: json.Marshal() fail: ", err)
			return
		}
		w.Write(jsonResp)
	}
}


func main() {
	flag.Parse()
	fhlog, err := os.OpenFile(*logFile, os.O_APPEND |os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatal(err)
	}

	defer fhlog.Close()
	log.SetOutput(fhlog)
	log.SetFlags(log.Lmicroseconds | log.LUTC | log.Ldate | log.Ltime)
	log.Println("logging initialized")
	log.Println("using", *conf)
	data, err := ioutil.ReadFile(*conf)
	if err != nil {
		fmt.Println(*conf,err)
		log.Fatal(*conf, err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(*conf, err)
		log.Fatal(*conf, err)
	}

	router := mux.NewRouter()
	router.Use(middleware)
	router.HandleFunc("/api/v1/", root)
	router.HandleFunc("/api/v1/echo", echo)
	router.HandleFunc("/api/v1/endpoint", endpoint)
	log.Println("binding on port ", config.Port)
	log.Println("rest_api starting")
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(config.Port), router))
}

