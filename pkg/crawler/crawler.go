package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"time"

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

func init() {
	RegisterCrawler("penti", &PentiCrawler{})
	RegisterCrawler("fuli", &FuliCrawler{})
}

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
	for name, crawler := range crawlers {
		err := json.Unmarshal(configs[name], crawler)
		if err != nil {
			log.Fatal(err)
		}
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
	go processArticle(db)
	for {
		cst := time.FixedZone("CST", 8*3600) // 东八
		now := time.Now().In(cst)
		if now.Hour() >= 16 && now.Hour() < 18 && now.Minute()%10 == 0 {
			for _, crawler := range crawlers {
				crawler.Crawl()
			}
		}
		time.Sleep(time.Minute)
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
				fmt.Printf("new article: %s\n", article.Title)
				db.Create(&article)
				articles = append(articles, article)
				dingTalkClient.PushArticles(articles)
			}
		} else {
			fmt.Printf("article exist: %s\n", article.Title)
		}
	}
}
