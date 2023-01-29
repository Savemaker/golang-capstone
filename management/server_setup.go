package management

import "github.com/labstack/echo/v4"

func LocationServiceServerSetup(l *LocationService) *echo.Echo {
	handlers := &handler{locationService: l}
	e := echo.New()
	e.POST("/location", handlers.PostLocation)
	e.GET("/users", handlers.GetUsers)
	return e
}
