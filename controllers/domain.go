package controllers

import (
	"github.com/theredwiking/domainscan/models"
	"sync"
)

type Domain struct {
	Url string `json:"domain" form:"domain"`
}

func (d Domain) Start() (models.Nmap, models.Headers, error) {
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
}
