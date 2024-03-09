package controllers

import (
	"database/sql"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/theredwiking/domainscan/models"
	"github.com/theredwiking/domainscan/views"
)

func HomeHandler(c echo.Context) error {
	return render(c, views.Home("Welcome"))
}

func ScanResultHandler(c echo.Context, scanned models.Nmap, get models.Headers) error {
	return render(c, views.ScanResult(scanned, get))
}

func ScansHandler(c echo.Context, db *sql.DB) error {
	results, err := DomainResults(db)
	if err != nil {
		return err
	}
	return render(c, views.Scans("Scans", results))
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
