package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	m "location/management"

	_ "github.com/lib/pq"
)

func main() {
	log.Print("starting management service")

	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	loc := os.Getenv("POSTGRES_LOCATION")
	dbName := os.Getenv("POSTGRES_DB_NAME")

	log.Print("...establishing connection with db...")

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pwd, loc, dbName)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("OK - db connection established!")
	}

	db.Exec("CREATE TABLE IF NOT EXISTS users (user_name varchar PRIMARY KEY, latitude float, longitude float);")

	locationService := m.LocationService{Repository: &m.UserRepo{DB: db}, UserFinder: &m.HaversineUserFinder{}}
	echoServer := m.LocationServiceServerSetup(&locationService)
	echoServer.Logger.Fatal(echoServer.Start(":8082"))
}
