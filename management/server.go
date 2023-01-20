package management

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
)

func LocationManagementServer(db *sql.DB) *echo.Echo {
	l := LocationService{Repository: &PostgresqlRepo{db}, UserFinder: &HaversineUserFinder{}}
	e := echo.New()
	e.POST("/location", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}

		l.UpdateUserLocation(u)
		return c.NoContent(200)
	})

	e.GET("/users", func(c echo.Context) error {
		latitude := c.QueryParam("latitude")
		longitude := c.QueryParam("longitude")
		radius := c.QueryParam("radius")

		if latitude == "" || longitude == "" || radius == "" {
			return c.String(400, "Provide latitude, longitude and radius as query params")
		}

		coordRegex := regexp.MustCompile(`^-?(90|[-8]?\d(\.\d+)?), -?(180|[01]?[-7]?\d(\.\d+)?)$`)
		lat, err := strconv.ParseFloat(latitude, 64)
		if err != nil {
			log.Fatal(err)
		}
		lon, err := strconv.ParseFloat(longitude, 64)
		if err != nil {
			log.Fatal(err)
		}
		rad, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			log.Fatal(err)
		}
		users := l.SearchUsersNearby(&Location{Latitude: lat, Longitude: lon}, rad)
		return c.JSON(200, users)
	})
	return e
}
