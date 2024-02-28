package main

import (
	"github.com/theredwiking/domainscan/routes"
	"net/http"
	"embed"
	"log"
	"os"
	"io/fs"

	"github.com/labstack/echo/v4"
)

//go:embed assets/*
var embededFiles embed.FS

func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("using live mode")
		return http.FS(os.DirFS("assets"))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(embededFiles, "assets")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func main() {
	e := echo.New()

	e.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is running")
	})

	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useOS))
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	routes.Views(e)
	routes.Domain(e)

	e.Logger.Fatal(e.Start(":3000"))
}
