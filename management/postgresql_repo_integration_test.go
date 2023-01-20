package management

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgresqlRepo(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("CREATE TABLE users (user_name varchar PRIMARY KEY, latitude float, longitude float);")
	db.Exec(`	INSERT INTO users (user_name, latitude, longitude) 
				VALUES 
						('stepa',  44.842465, 20.380679), 
						('john', 51.507322, -0.127647);`)

	repo := PostgresqlRepo{DB: db}

	for desc, testFunc := range map[string]func(*testing.T, *PostgresqlRepo){
		"test find all users":       FindAll,
		"test update user location": UpdateUserLocation,
	} {
		t.Run(desc, func(t *testing.T) {
			testFunc(t, &repo)
		})
	}
}

func FindAll(t *testing.T, repo *PostgresqlRepo) {
	users := repo.FindAll()
	if len(users) == 0 {
		t.Error("error")
	}
}

func UpdateUserLocation(t *testing.T, repo *PostgresqlRepo) {
	repo.UpdateUserLocation(&User{Name: "stepa", UserLocation: &Location{Latitude: 69, Longitude: 420}})
	repo.UpdateUserLocation(&User{Name: "weirdo", UserLocation: &Location{Latitude: -69, Longitude: -420}})
	users := repo.FindAll()
	if len(users) != 3 {
		t.Error("error")
	}
	var userName string
	var latitude, longitude float64
	repo.DB.QueryRow("SELECT * FROM users WHERE user_name='stepa';").Scan(&userName, &latitude, &longitude)
	if latitude != 69 || longitude != 420 {
		t.Error("fucked up")
	}
}
