package controllers

import (
	"context"
	"github.com/theredwiking/domainscan/models"
	"log"
	"sync"
	"time"

	"github.com/Ullaakut/nmap/v3"
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
