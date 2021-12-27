## DATABASE
---

#### Make the database network accessible
To make the new DB accessible on the network.  Edit /some/path/postgresql.conf and make the following changes

    #------------------------------------------------------------------------------
    # CONNECTIONS AND AUTHENTICATION
    #------------------------------------------------------------------------------

    # - Connection Settings -

    listen_addresses = '*'                                                   <= ADD
    #listen_addresses = 'localhost'		# what IP address(es) to listen on;
					# comma-separated list of addresses;
					# defaults to 'localhost'; use '*' for all
					# (change requires restart)

#### Provide user access to the database
Give the new admin and user access to the db.  Edit /some/path/pg_hba.conf

    # Database administrative login by Unix domain socket
    local   all             postgres                                peer

    # TYPE  DATABASE        USER            ADDRESS                 METHOD

    # "local" is for Unix domain socket connections only
    local   mydb_dev         all                                     md5      <= ADD LINE
    local   all             all                                     peer
    # IPv4 local connections:
    host    mynewdb        newdbuser         0.0.0.0/0              md5      <= ADD LINE
    host    all             all             127.0.0.1/32            md5


#### Restart the database

When the edits are complete, restart PostgreSQL

    > systemctl restart postgresql
    > systemctl status postgresql
    ● postgresql.service - PostgreSQL RDBMS
         Loaded: loaded (/lib/systemd/system/postgresql.service; enabled; vendor preset: enabled)
         Active: active (exited) since Sun 2021-12-26 09:41:53 MST; 4s ago
        Process: 35589 ExecStart=/bin/true (code=exited, status=0/SUCCESS)
        Main PID: 35589 (code=exited, status=0/SUCCESS)

        Dec 26 09:41:53 coffee systemd[1]: Starting PostgreSQL RDBMS...
        Dec 26 09:41:53 coffee systemd[1]: Finished PostgreSQL RDBMS.

#### Execute script/db_setup

    > db_setup
    Database to be used in
    	[P]roduction
	    [Q]A/Test
	    [D]evelopment

            [P/q/d] : d

       Database name: mydb
          Admin user: db_admin
	      Set password for "db_admin":
          Confirm password:
    Application user: db_user
          Set password for "db_user":
          Confirm password:

       DB_NAME: mydb_dev
     DB_SCHEMA: mydb_dev
      DB_ADMIN: db_admin
       DB_USER: db_user


    Create script [Y/n]?

    Creating mydb_dev.sql

    Run "psql -f mydb_dev.sql" to create database
    >

#### Execute the generated SQL

    > sql -f mydb_dev.sql
    CREATE ROLE
    CREATE ROLE
    CREATE DATABASE
    You are now connected to database "mydb_dev" as user "postgres".
    DROP SCHEMA
    CREATE SCHEMA
    ALTER ROLE
    ALTER ROLE
    Password for user db_admin:
    You are now connected to database "mydb_dev" as user "db_admin".
    GRANT
    ALTER DEFAULT PRIVILEGES
    ALTER DEFAULT PRIVILEGES
    ALTER DEFAULT PRIVILEGES
    ALTER DEFAULT PRIVILEGES
   >

#### Check the new database

    postgres=# \l
                                  List of databases
    Name    |  Owner   | Encoding |   Collate   |    Ctype    |   Access privileges
    -----------+----------+----------+-------------+-------------+-----------------------
    mydb_dev  | db_admin | UTF8     | en_US.UTF-8 | en_US.UTF-8 |

    postgres=# \du
                                   List of roles
    Role name |                         Attributes                         | Member of
    -----------+------------------------------------------------------------+-----------
    db_admin  |                                                            | {}
    db_user   |                                                            | {}

Connect as db_admin and test a table create/drop

    postgres=# \c mydb_dev db_admin
    Password for user db_admin:
    You are now connected to database "mydb_dev" as user "db_admin".
    mydb_dev=> create table test();
    CREATE TABLE
    mydb_dev=> \d
         List of relations
    Schema  | Name | Type  |  Owner
    ----------+------+-------+----------
    mydb_dev | test | table | db_admin

    mydb_dev=> drop table test;
    DROP TABLE

Try the same with db_user

    mydb_dev=> \c mydb_dev db_user;
    Password for user db_user:
    You are now connected to database "mydb_dev" as user "db_user".
    mydb_dev=> create table test();
    ERROR:  permission denied for schema mydb_dev                  <= WHAT WE WANT
    LINE 1: create table test();
                          ^
     mydb_dev=> \q

Last, test db_user access via network

    > psql -U db_user -h localhost mydb_dev
    Password for user db_user:
    psql (12.9 (Ubuntu 12.9-0ubuntu0.20.04.1))
    SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
    Type "help" for help.

    mydb_dev=>

#### Create a table for testing

    CREATE TABLE test (
    time     timestamp,
    key      varchar(32) not null,
    field1   varchar(32),
    field2   varchar(32)
    );


