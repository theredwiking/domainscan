package controllers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/theredwiking/domainscan/models"

	"github.com/Ullaakut/nmap/v3"
	parser "github.com/blockloop/scan"
)

func scan(domain string) (models.Nmap, error) {
	var complete models.Nmap
	complete.Url = domain
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(domain),
	)

	if err != nil {
		return models.Nmap{}, err
	}

	result, warnings, err := scanner.Run()
	if len(*warnings) > 0 {
		log.Printf("finished run with warnings: %s\n", *warnings)
	}
	if err != nil {
		return models.Nmap{}, err
	}

	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		complete.Ip = host.Addresses[0].String()
		var ports []models.Port
		
		for _, port := range host.Ports {
			temp := models.Port{port.ID, port.Protocol, port.State.String(), port.Service.Name}
			ports = append(ports, temp)
		}

		complete.Ports = ports
	}
	return complete, nil
}

func runNmap(wg *sync.WaitGroup, channel chan models.Nmap, domain string) {
	defer wg.Done()
	result, err := scan(domain)
	if err != nil {
		log.Printf("Error with headers: %v", err)
		channel <- models.Nmap{}
	}
	channel <- result
}

func insertPortsDB(db *sql.DB, scan models.Nmap, domainId int64) {
	query, err := db.Prepare("INSERT INTO nmap (port, protocol, service) VALUES (?, ?, ?)")		
	if err != nil {
		log.Printf("Failed to prepare port statement: %v", err)
	}

	connectionQuery, err := db.Prepare("INSERT INTO domainNmap (domainId, nmapId, state) VALUES (?, ?, ?)")		
	if err != nil {
		log.Printf("Failed to prepare domainNmap statement: %v", err)
	}

	go insertIpDB(db, domainId, scan.Ip)

	for _, port := range scan.Ports  {
		var temp int
		err := db.QueryRow("SELECT id FROM nmap WHERE port = ?", port.Port).Scan(&temp)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			result , err := query.Exec(port.Port, port.Protocol, port.Service)
			if err != nil {
				log.Printf("Failed to insert port: %v", err)
			}

			id, err := result.LastInsertId()
			if err != nil {
				log.Printf("Failed to get id: %v", err)
				return
			}
			if id != 0 {
				_, err = connectionQuery.Exec(domainId, id, port.State)
				if err != nil {
					log.Printf("Failed to insert domainNmap: %v", err)
				}
			}
		} else if err != nil {
			log.Printf("Failed to check if port exist: %v", err)
		}
		if temp != 0 {

			_, err = connectionQuery.Exec(domainId, temp, port.State)
			if err != nil {
				log.Printf("Failed to insert domainNmap: %v", err)
			}
		}
	}
	return
}

func fromDb(db *sql.DB, id int64) models.Nmap {
	var result models.Nmap
	rows, err := db.Query("SELECT nmapId, state FROM domainNmap WHERE domainId = ?", id)
	if err != nil {
		return models.Nmap{}
	}

	for rows.Next() {
		var id int
		var port models.Port
		if err := rows.Scan(&id, &port.State); err != nil {
			log.Printf("Error parse: %v", err)
		}

		row, err := db.Query("SELECT port, protocol, service FROM nmap WHERE id = ?", id)
		if err != nil {
			log.Printf("Failed to find row: %v", err)
		}
		if err := parser.Row(&port, row); err != nil {
			log.Printf("failed to parse port: %v", err)
		}
		result.Ports = append(result.Ports, port)
	}
	return result
}
