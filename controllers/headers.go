package controllers

import (
	"github.com/theredwiking/domainscan/models"
	"fmt"
	"log"
	"net/http"
	"sync"
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
