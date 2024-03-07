package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/theredwiking/domainscan/models"
)

func headers(domain string) (models.Headers, error) {
	var result models.Headers

	resp, err := http.Get(fmt.Sprintf("https://%s", domain))
	if err != nil {
		return models.Headers{}, err
	}

	result.Protocol = resp.Proto
	result.ContentType = resp.Header["Content-Type"][0]
	result.Server = resp.Header["Server"][0]
	
	return result, nil	
	/*body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}*/
}

func runHeaders(wg *sync.WaitGroup, channel chan models.Headers, domain string) {
	defer wg.Done()
	result, err := headers(domain)
	if err != nil {
		log.Printf("Error with headers: %v", err)
		channel <- models.Headers{}
	}
	channel <- result
}

func insertDomainDB(db *sql.DB, name string, headers models.Headers) int64 {
	query, err := db.Prepare("INSERT INTO domain (name, protocol, contentType, server) VALUES (?, ?, ?, ?)")		
	if err != nil {
		log.Printf("Failed to prepare domain statement: %v", err)
		return 0
	}
	result, err := query.Exec(name, headers.Protocol, headers.ContentType, headers.Server)
	if err != nil {
		log.Printf("Failed to insert domain: %v", err)
		return 0
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get id: %v", err)
		return 0
	}
	return id
}

func insertIpDB(db *sql.DB, domainId int64, ip string) {
	if domainId != 0 {
		query, err := db.Prepare("UPDATE domain SET ip = ? WHERE id = ?")
		if err != nil {
			log.Printf("Failed to insert ip: %v", err)
			return
		}
		_, err = query.Exec(ip, domainId)
		if err != nil {
			log.Printf("Failed to run ip query: %v", err)
			return
		}
	}
	return
}
