package controllers

import (
	"database/sql"
	"log"
	"sync"

	"github.com/theredwiking/domainscan/models"
)

type Domain struct {
	Url string `json:"domain" form:"domain"`
}

func (d Domain) Start(db * sql.DB) (models.Nmap, models.Headers, error) {
	var headers models.Headers
	rows, err := db.Query("SELECT id, ip, protocol, contentType, server FROM domain WHERE name = ?", d.Url)
	if err != nil {
		log.Printf("Failed to check domain: %v", err)
		return models.Nmap{}, models.Headers{}, err
	}
	for rows.Next() {
		var ip string
		if err := rows.Scan(&headers.Id, &ip, &headers.Protocol, &headers.ContentType, &headers.Server); err != nil {
			log.Printf("Failed to parse domain: %v", err)
			return models.Nmap{}, models.Headers{}, err
		}

		scan := fromDb(db, headers.Id)
		scan.Ip = ip

		return scan, headers, nil
	}
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

	go insertingDB(db, d.Url, get, scan)

	return scan, get, nil

}

func DomainResults(db *sql.DB) ([]models.Combined ,error) {
	var results []models.Combined

	rows, err := db.Query("SELECT id, ip, protocol, contentType, server FROM domain")
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var ip string
		var header models.Headers

		if err := rows.Scan(&header.Id, &ip, &header.Protocol, &header.ContentType, &header.Server); err != nil {
			return results, err
		}
		scan := fromDb(db, header.Id)
		scan.Ip = ip
		combined := models.Combined{Headers: header, Scans: scan}
		results = append(results, combined)
	}

	return results, nil
}

func insertingDB(db *sql.DB, domain string, headers models.Headers, scan models.Nmap) {
	domainId := insertDomainDB(db, domain, headers)
	insertPortsDB(db, scan, domainId)
}
