package main

import (
	sql "github.com/theredwiking/domainscan/database"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/theredwiking/domainscan/routes"

	"github.com/labstack/echo/v4"
)

//go:embed schema.sql
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

	db, err := sql.Connect()
	if err != nil {
		log.Fatalf("Error with database: %v\n", err)
	}

	e.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is running")
	})

	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useOS))
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	routes.Views(e)
	routes.Domain(e, db)

	e.Logger.Fatal(e.Start(":3000"))
}
