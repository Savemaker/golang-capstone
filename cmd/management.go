package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"location/management"

	_ "github.com/lib/pq"
)

func main() {
	log.Print("starting management service")

	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	loc := os.Getenv("POSTGRES_LOCATION")

	log.Print("establishing connection with db")

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s/postgres?sslmode=disable", user, pwd, loc)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("connection established")
	}

	echoServer := management.LocationManagementServer(db)
	echoServer.Logger.Fatal(echoServer.Start(":8082"))
}
