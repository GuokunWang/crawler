package crawler

import (
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
)

type KexueCrawler struct {
	Url string `json:"url"`
}

func (k *KexueCrawler) Crawl() error {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(k.Url)
	for _, item := range feed.Items {
		article := Article{
			ID:    uuid.New().String(),
			Title: item.Title,
			URL:   item.Link,
		}
		articleChan <- article
	}
	return nil
}
