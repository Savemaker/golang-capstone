package management

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgresqlRepo(t *testing.T) {
	if f := os.Getenv("INTEGRATION_TESTS"); f == "" {
		log.Print("set INTEGRATION_TESTS env var to run integration test TestPostgresqlRepo")
		t.Skip()
	} else {
		log.Print("integration tests are enabled")
	}
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/location_test?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("CREATE TABLE users (user_name varchar PRIMARY KEY, latitude float, longitude float);")
	db.Exec(`	INSERT INTO users (user_name, latitude, longitude) 
				VALUES 
						('stepa',  44.842465, 20.380679), 
						('john', 51.507322, -0.127647);`)

	repo := UserRepo{DB: db}

	for desc, testFunc := range map[string]func(*testing.T, *UserRepo){
		"test find all users":       FindAll,
		"test update user location": UpdateUserLocation,
	} {
		t.Run(desc, func(t *testing.T) {
			testFunc(t, &repo)
		})
	}
	db.Exec("DROP TABLE users")
}

func FindAll(t *testing.T, repo *UserRepo) {
	users, _ := repo.FindAll()
	if len(users) == 0 {
		t.Error("expected more than 0 users")
	}
}

func UpdateUserLocation(t *testing.T, repo *UserRepo) {
	repo.UpdateUserLocation(&User{Name: "stepa", UserLocation: &Location{Latitude: 69, Longitude: 420}})
	repo.UpdateUserLocation(&User{Name: "weirdo", UserLocation: &Location{Latitude: -69, Longitude: -420}})
	users, _ := repo.FindAll()
	if len(users) != 3 {
		t.Error("expected 3 users")
	}
	var userName string
	var latitude, longitude float64
	repo.DB.QueryRow("SELECT * FROM users WHERE user_name='stepa';").Scan(&userName, &latitude, &longitude)
	if latitude != 69 || longitude != 420 {
		t.Error("failed to update user coordinates")
	}
}
