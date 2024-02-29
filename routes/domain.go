package routes

import (
	"database/sql"
	"log"

	"github.com/theredwiking/domainscan/controllers"

	"github.com/labstack/echo/v4"
)

func Domain(e *echo.Echo, db *sql.DB) {
	api := e.Group("/api")

	api.POST("/scan", func(c echo.Context) error {
		domain := new(controllers.Domain)
		if err := c.Bind(domain); err != nil {
			log.Fatalf("Failed to parse: %v", err)
			return err
		}

		scanned, headers, err:= domain.Start(db);
		if err != nil {
			return err
		}

		return controllers.ScanResultHandler(c, scanned, headers)
	})
}
