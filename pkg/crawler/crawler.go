package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Crawler interface{ Crawl() error }

type Article struct {
	ID    string
	Title string `gorm:"primaryKey"`
	URL   string
}

var crawlers = make(map[string]Crawler)

var articleChan = make(chan Article)

var dingTalkClient = DingTalkClient{}

func ConfigCrawler(confPath string) {
	// read config file
	confFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal(err)
	}

	configs := make(map[string]json.RawMessage)
	err = json.Unmarshal(confFile, &configs)
	if err != nil {
		log.Fatal(err)
	}

	// init dingTalkClient first
	err = json.Unmarshal(configs["dingTalk"], &dingTalkClient)
	if err != nil {
		log.Fatal(err)
	}

	// init other crawler in map
	rssCrawlers := []RssCrawler{}
	err = json.Unmarshal(configs["rssCrawlers"], &rssCrawlers)
	if err != nil {
		log.Fatal(err)
	}

	for i, crawler := range rssCrawlers {
		RegisterCrawler(crawler.Name, &rssCrawlers[i])
	}
}

func RegisterCrawler(name string, crawler Crawler) {
	v := reflect.TypeOf(crawler)
	if v == nil || v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("client [%s] muse be registered with pointer, but got %s", name, v.Kind()))
	}
	if _, ok := crawlers[name]; ok {
		panic(fmt.Sprintf("client [%s] already exists", name))
	}
	crawlers[name] = crawler
}

func Run(db *gorm.DB) {
	r := rand.New(rand.NewSource(55))
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	interval := 5 + r.Intn(6)

	go processArticle(db)
	for {
		log.Infof("sleep for %d minute", interval)
		time.Sleep(time.Duration(interval) * time.Minute)
		for _, crawler := range crawlers {
			crawler.Crawl()
		}
		interval = 5 + r.Intn(6)
	}
}

func processArticle(db *gorm.DB) {
	for {
		articles := []Article{}
		article, ok := <-articleChan
		if !ok {
			break
		}
		if err := db.First(&article, "url = ?", article.URL).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Infof("new article: %s\n", article.Title)
				db.Create(&article)
				articles = append(articles, article)
				dingTalkClient.PushArticles(articles)
			}
		} else {
			log.Infof("article exist: %s\n", article.Title)
		}
	}
}
