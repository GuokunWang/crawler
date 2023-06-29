package crawler

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type PentiCrawler struct {
	Url    string `json:"url"`
	Perfix string `json:"prefix"`
	Key    string `json:"key"`
}

func (c *PentiCrawler) Crawl() error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(c.Url)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer resp.Body.Close()

	reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if strings.Contains(string(body), c.Key) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			log.Fatalln(err)
			return err
		}
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			title := s.Text()
			if strings.Contains(title, c.Key) {
				href, _ := s.Attr("href")
				title, _ := s.Attr("title")
				if len(title) > 3 {
					article := Article{
						ID:    uuid.New().String(),
						Title: title,
						URL:   c.Perfix + href,
					}
					articleChan <- article
				}
			}
		})
	}
	return nil
}
