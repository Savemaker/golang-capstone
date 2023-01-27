package management

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

func LocationManagementServer(db *sql.DB) *echo.Echo {
	l := LocationService{Repository: &UserRepo{db}, UserFinder: &HaversineUserFinder{}}
	e := echo.New()
	e.POST("/location", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		if !ValidateUserName(u.Name) {
			return c.String(400, "username: 4-16 symbols (a-zA-Z0-9 symbols are acceptable)")
		}
		if _, err := ValidateLatitudeAndGet(fmt.Sprintf("%f", u.UserLocation.Latitude)); err != nil {
			return c.String(400, "latitude: "+err.Error())
		}
		if _, err := ValidateLongitudeAndGet(fmt.Sprintf("%f", u.UserLocation.Longitude)); err != nil {
			return c.String(400, "longitude: "+err.Error())
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
		lat, err := ValidateLatitudeAndGet(latitude)
		if err != nil {
			return c.String(400, "latitude: "+err.Error())
		}
		lon, err := ValidateLongitudeAndGet(longitude)
		if err != nil {
			return c.String(400, "longitude: "+err.Error())
		}
		rad, err := ValidateRadiusAndGet(radius)
		if err != nil {
			return c.String(400, "radius: "+err.Error())
		}
		users := l.SearchUsersNearby(&Location{Latitude: lat, Longitude: lon}, rad)
		return c.JSON(200, users)
	})
	return e
}
