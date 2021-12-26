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
	Host	string
	Port	int
	User	string
	Passwd  string
	DBName  string
	Sslmode	string
}

var (
	port	= flag.Int
	conf	= flag.String
	logFile =
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w, http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
}

func apiRoot(w, http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp := make(map[string]string)
	resp["status"] = "ok"
	resp["message"] = "success"
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func echo(w, http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp := make(map[string]string)
	resp["status"] = "ok"
	resp["message"] = "echo"
	jsonResp, err := json.Marshal(resp)
	W.Write(jsonResp)
}

func endpoint(w, http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp := make(map[string]string)
	var data InboundPayload

	if r.Method == "POST" {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http:StatusInternalServerError)
			resp["status"] = "not_ok"
			resp["message"] = err.Error()
			log.Println(resp["message"])
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w., "{\"status\":\"not_ok\",\"message\":\"json.Marshal() fail\"}")
				log.Println("json.Marshal() fail", err)
				return
			}
			W.Write(jsonResp)
			return
	}

		err = json.Unmarshal(reqBody, &data)
		if err != nil {
			w.WriteHeader(http:StatusInternalServerError)
			resp["status"] = "not_ok"
			resp["message"] = err.Error()
			log.Println(resp["message"])
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w., "{\"status\":\"not_ok\",\"message\":\"json.Marshal() fail\"}")
				log.Println("json.Marshal() fail", err)
				return
			}
			W.Write(jsonResp)
			return
		}

		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
					config.Host, config.Port, config.User, config.Password, config.Dbname, config.Sslmode)

		conn,err := sql.Open("postgres", connString)
		if err != nil {
			w.WriteHeader(http:StatusInternalServerError)
			resp["status"] = "not_ok"
			resp["message"] = err.Error()
			log.Println(resp["message"])
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w., "{\"status\":\"not_ok\",\"message\":\"json.Marshal() fail\"}")
				log.Println("json.Marshal() fail", err)
				return
			}
			W.Write(jsonResp)
			return
		}

		defer conn.Close()
		statement := fmt.Sprintf("insert into app_table (time, key, field1, field2) values ('%s', '%s', '%s', '%s');",
					data.Time, data.Key, data.Field1, data.Field2)

		_, err = conn.Exec(statement)
		if err != nil {
			w.WriteHeader(http:StatusInternalServerError)
			resp["status"] = "not_ok"
			resp["message"] = err.Error()
			log.Println(resp["message"])
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Fprintf(w, "{\"status\":\"not_ok\",\"message\":\"json.Marshal() fail\"}")
				log.Println("conn.Exec(): json.Marshal() fail: ", err)
				return
			}
			w.Write(jsonResp)
			return
		}
	} else if r.Method == "GET" {
		w.WriteHeader(http:StatusNotFound)
		resp["status"] = "ok"
		resp["message"] = "method not implemented"
		log.Println(resp["message"])
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "{\"status\":\"not_ok\",\"message\":\"json.Marshal() fail\"}")
			log.Println("conn.Exec(): json.Marshal() fail: ", err)
			return
		{
		w.Write(jsonResp)
	} else {
		w.WriteHeader(http:StatusMethodNotAllowed))
		resp["status"] = "ok"
		resp["message"] = "method not allowed"
		log.Println(resp["message"])
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
	if err != nid {
		log.Fatal(err)
	}

	defer fhlog.Close()
	log.SetOutput(fhlog)
	log.SetFlags(log.Lmicroseconds | log.LUTC | log.Ldate | log.Ltime)
	log.Println("logging initialized")
	log.Println("using", *conf)
	data, err := ioutil.ReadFile(*conf)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &config)
		fmt.Println(err)
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.Use(middleware)
	router.HandleFunc("api/v1/", apiRoot)
	router.HandleFunc("api/v1/echo", echo)
	router.HandleFunc("api/v1/endpoint", endpoint)
	log.Println("binding on port ", *port)
	log.Println("rest_api starting")
	log.Fatal(http.ListenAndServe(":" + strconf.Itoa(*port), router)
}

