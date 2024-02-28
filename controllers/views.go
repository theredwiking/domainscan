package controllers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/theredwiking/domainscan/models"
	"github.com/theredwiking/domainscan/views"
)

func HomeHandler(c echo.Context) error {
	return render(c, views.Home("Welcome"))
}

func ScanResultHandler(c echo.Context, result models.Result) error {
	return render(c, views.ScanResult(result))
}

func ScansHandler(c echo.Context) error {
	return render(c, views.Scans("Scans"))
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
