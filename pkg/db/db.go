package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const connStrCfg = "%s:%s@tcp(%s:%s)/guild_cfg"
const connStrScraper = "%s:%s@tcp(%s:%s)/scraper"

var DBCfg *sql.DB
var DBScraper *sql.DB

func IsDuplicateErr(err error) bool {
	return strings.Contains(err.Error(), "Error 1062 (23000): Duplicate entry")
}

func initDB(connStr string) *sql.DB {
	var ret *sql.DB
	// connect to DB and set objects
	user := os.Getenv("db_user")
	pass := os.Getenv("db_pass")
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")

	ret, err := sql.Open("mysql", fmt.Sprintf(connStr, user, pass, host, port))
    if err != nil {
        log.Fatalf("error on open: %v", err)
    }
	err = ret.Ping()
    if err != nil {
        log.Fatalf("error on connect: %v", err)
    }
	return ret
}

func init() {
	// read env
	err := godotenv.Load()
	if err != nil {
		log.Printf("WARNING: could not load .env: %v", err)
	}

	log.Println("Initializing server config db...")
	DBCfg = initDB(connStrCfg)

	log.Println("Initializing scraper db...")
	DBScraper = initDB(connStrScraper)

}