package routes

import (
	"github.com/theredwiking/domainscan/controllers"

	"github.com/labstack/echo/v4"
)

func Views(e *echo.Echo) {
	group := e.Group("/")
	group.GET("", controllers.HomeHandler)
	group.GET("scans", controllers.ScansHandler)
}
