package routes

import (
	"database/sql"

	"github.com/theredwiking/domainscan/controllers"

	"github.com/labstack/echo/v4"
)

func Views(e *echo.Echo, db *sql.DB) {
	group := e.Group("/")
	group.GET("", controllers.HomeHandler)
	group.GET("scans", func (c echo.Context) error {return controllers.ScansHandler(c, db)})
}
