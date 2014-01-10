package odb

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func fetch(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unsuccessful status code %d", res.StatusCode)
	}
	q, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	q.Url = res.Request.URL
	return q, nil
}
