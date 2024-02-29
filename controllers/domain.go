package controllers

import (
	"database/sql"
	"log"
	"sync"

	parser "github.com/blockloop/scan"

	"github.com/theredwiking/domainscan/models"
)

type Domain struct {
	Url string `json:"domain" form:"domain"`
}

func (d Domain) Start(db * sql.DB) (models.Nmap, models.Headers, error) {
	var headers models.Headers
	rows, err := db.Query("SELECT id, ip, protocol, server FROM domain WHERE name = ?", d.Url)
	if err != nil {
		return models.Nmap{}, models.Headers{}, err
	}

	if !rows.Next() {
		wg := new(sync.WaitGroup)
		wg.Add(2)

		header := make(chan models.Headers)
		scanned := make(chan models.Nmap)
		defer close(header)
		defer close(scanned)

		go runHeaders(wg, header, d.Url)
		go runNmap(wg, scanned, d.Url)

		get := <-header
		scan := <-scanned
		wg.Wait()

		return scan, get, nil
	} else {
		if err := parser.Row(&headers, rows); err != nil {
			log.Printf("Failed to parse: %v", err)
		}

		scan := fromDb(db, headers.Id)

		return scan, headers, nil
	}
}
