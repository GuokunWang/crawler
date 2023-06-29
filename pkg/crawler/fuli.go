package crawler

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

type FuliCrawler struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

func (c *FuliCrawler) Crawl() error {
	doc, err := goquery.NewDocument(c.Url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), c.Key) {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				url, _ := s.Attr("href")
				title := s.Text()
				article := Article{
					ID:    uuid.New().String(),
					Title: title,
					URL:   url,
				}
				articleChan <- article
			})
		}
	})
	return nil
}
