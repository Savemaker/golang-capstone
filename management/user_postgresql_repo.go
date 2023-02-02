package management

import (
	"database/sql"
	"errors"
	"log"
)

type UserRepo struct {
	DB *sql.DB
}

var (
	ErrInternalServerError = errors.New("internal server error")
)

func (p *UserRepo) FindAll() ([]User, error) {
	users := make([]User, 0)

	s, err := p.DB.Prepare("SELECT * FROM users")
	if err != nil {
		log.Print(err)
		return users, ErrInternalServerError
	}
	res, err := s.Query()
	if err != nil {
		log.Print(err)
		return users, ErrInternalServerError
	}
	defer res.Close()
	for res.Next() {
		var userName string
		var latitude float64
		var longitude float64

		err := res.Scan(&userName, &latitude, &longitude)
		if err != nil {
			log.Print(err)
			return users, ErrInternalServerError
		}
		users = append(users, User{Name: userName, UserLocation: &Location{Latitude: latitude, Longitude: longitude}})
		if err := res.Err(); err != nil {
			log.Print(err)
			return users, ErrInternalServerError
		}
	}
	return users, nil
}

func (p *UserRepo) UpdateUserLocation(user *User) error {
	st, err := p.DB.Prepare("SELECT * FROM users WHERE user_name=$1;")

	if err != nil {
		log.Print(err)
		return ErrInternalServerError
	}

	var userName string
	var latitude float64
	var longitude float64

	err = st.QueryRow(user.Name).Scan(&userName, &latitude, &longitude)

	if err == sql.ErrNoRows {
		st, err := p.DB.Prepare("INSERT INTO users (user_name, latitude, longitude) VALUES ($1, $2, $3);")

		if err != nil {
			log.Print(err)
			return ErrInternalServerError
		}
		_, err = st.Exec(user.Name, user.UserLocation.Latitude, user.UserLocation.Longitude)
		if err != nil {
			log.Print(err)
			return ErrInternalServerError
		}
	} else {
		st, err := p.DB.Prepare("UPDATE users SET latitude=$1, longitude=$2 WHERE user_name=$3;")

		if err != nil {
			log.Print(err)
			return ErrInternalServerError
		}
		_, err = st.Exec(user.UserLocation.Latitude, user.UserLocation.Longitude, user.Name)
		if err != nil {
			log.Print(err)
			return ErrInternalServerError
		}
	}
	return nil
}
