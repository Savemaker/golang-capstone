package management

import (
	"database/sql"
	"log"
)

type PostgresqlRepo struct {
	DB *sql.DB
}

func (p *PostgresqlRepo) FindAll() []User {
	users := make([]User, 0)

	s, err := p.DB.Prepare("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	res, err := s.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()
	for res.Next() {
		var userName string
		var latitude float64
		var longitude float64

		err := res.Scan(&userName, &latitude, &longitude)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, User{Name: userName, UserLocation: &Location{Latitude: latitude, Longitude: longitude}})
		if err := res.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return users
}

func (p *PostgresqlRepo) UpdateUserLocation(user *User) {
	st, err := p.DB.Prepare("SELECT * FROM users WHERE user_name=$1;")

	if err != nil {
		log.Fatal(err)
	}

	var userName string
	var latitude float64
	var longitude float64

	err = st.QueryRow(user.Name).Scan(&userName, &latitude, &longitude)

	if err == sql.ErrNoRows {
		st, err := p.DB.Prepare("INSERT INTO users (user_name, latitude, longitude) VALUES ($1, $2, $3);")

		if err != nil {
			log.Fatal(err)
		}
		_, err = st.Exec(user.Name, user.UserLocation.Latitude, user.UserLocation.Longitude)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		st, err := p.DB.Prepare("UPDATE users SET latitude=$1, longitude=$2 WHERE user_name=$3;")

		if err != nil {
			log.Fatal(err)
		}
		_, err = st.Exec(user.UserLocation.Latitude, user.UserLocation.Longitude, user.Name)
		if err != nil {
			log.Fatal(err)
		}
	}
}
