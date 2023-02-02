package management

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

type Response struct {
	Users []User `json:"users"`
}

type handler struct {
	locationService *LocationService
}

func (h *handler) PostLocation(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if err := ValidateUserName(u.Name); err != nil {
		return c.JSON(400, ErrorResponse{Reason: err.Error()})
	}
	if _, err := ValidateLatitudeAndGet(fmt.Sprintf("%f", u.UserLocation.Latitude)); err != nil {
		return c.JSON(400, ErrorResponse{Reason: "latitude: " + err.Error()})
	}
	if _, err := ValidateLongitudeAndGet(fmt.Sprintf("%f", u.UserLocation.Longitude)); err != nil {
		return c.JSON(400, ErrorResponse{Reason: "longitude: " + err.Error()})
	}
	err := h.locationService.UpdateUserLocation(u)
	if err == ErrInternalServerError {
		return c.JSON(500, err.Error())
	}
	return c.NoContent(200)
}

func (h *handler) GetUsers(c echo.Context) error {
	latitude := c.QueryParam("latitude")
	longitude := c.QueryParam("longitude")
	radius := c.QueryParam("radius")

	if err := ValidateRequestParams(latitude, longitude, radius); err != nil {
		return c.JSON(400, ErrorResponse{Reason: err.Error()})
	}
	lat, err := ValidateLatitudeAndGet(latitude)
	if err != nil {
		return c.JSON(400, ErrorResponse{Reason: err.Error()})
	}
	lon, err := ValidateLongitudeAndGet(longitude)
	if err != nil {
		return c.JSON(400, ErrorResponse{Reason: err.Error()})
	}
	rad, err := ValidateRadiusAndGet(radius)
	if err != nil {
		return c.JSON(400, ErrorResponse{Reason: err.Error()})
	}
	users, err := h.locationService.SearchUsersNearby(&Location{Latitude: lat, Longitude: lon}, rad)
	if err == ErrInternalServerError {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, Response{Users: users})
}
