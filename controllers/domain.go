package controllers

import (
	"github.com/theredwiking/domainscan/models"
	"fmt"
	"sync"
)

type Domain struct {
	Url string `json:"domain" form:"domain"`
}

func (d Domain) Start() (models.Result, error) {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	header := make(chan models.Headers)
	scanned := make(chan models.Result)
	defer close(header)
	defer close(scanned)

	go runHeaders(wg, header, d.Url)
	go runNmap(wg, scanned, d.Url)

	fmt.Println(<-header)
	result := <-scanned
	wg.Wait()

	return result, nil
}
