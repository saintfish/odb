package odb

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func fetch(c *http.Client, url string) (*goquery.Document, error) {
	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// Odb served 404 pages for upcoming month but the document list may be there
	if res.StatusCode != 200 && res.StatusCode != 404 {
		return nil, fmt.Errorf("Unsuccessful status code %d", res.StatusCode)
	}
	q, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	q.Url = res.Request.URL
	return q, nil
}
