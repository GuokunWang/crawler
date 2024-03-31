package crawler

import (
	"strings"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

type RssCrawler struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Key  string `json:"key"`
}

func (r *RssCrawler) Crawl() error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(r.Url)
	if err != nil {
		logrus.Errorf("init crawler [%s] with url [%s] failed: %s", r.Name, r.Url, err.Error())
		return err
	}
	if feed.Items == nil {
		logrus.Errorf("crawler [%s] with url [%s] have no feed", r.Name, r.Url)
		return nil
	}
	for _, item := range feed.Items {
		if len(r.Key) > 0 && !strings.Contains(item.Title, r.Key) {
			continue
		}
		article := Article{
			ID:    uuid.New().String(),
			Title: item.Title,
			URL:   item.Link,
		}
		articleChan <- article
	}
	return nil
}
