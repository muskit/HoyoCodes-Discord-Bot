package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const connStrCfg = "%s:%s@tcp(%s:%s)/server_cfg"
const connStrCodes = "%s:%s@tcp(%s:%s)/codes"

var DBCfg *sql.DB
var DBCodes *sql.DB

func init() {
	println("initializing db...")
	
	// read env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not load .env: %s", err)
	}

	// connect to DB and set objects
	user := os.Getenv("db_user")
	pass := os.Getenv("db_pass")
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")

	log.Println("Initializing server config db...")
	DBCfg, err = sql.Open("mysql", fmt.Sprintf(connStrCfg, user, pass, host, port))
    if err != nil {
        log.Fatalf("error on open: %v", err)
    }
	err = DBCfg.Ping()
    if err != nil {
        log.Fatalf("error on connect: %v", err)
    }
	log.Println("Server config db initialized!")

	// log.Println("Initializing codes db...")
	// DBCodes, err = sql.Open("mysql", fmt.Sprintf(connStrCodes, user, pass, host, port))
    // if err != nil {
    //     log.Fatalf("error on open: %v", err)
    // }
	// err = DBCodes.Ping()
    // if err != nil {
    //     log.Fatalf("error on connect: %v", err)
    // }
	// log.Println("Codes db initialized!")
}