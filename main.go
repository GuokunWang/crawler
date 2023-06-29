package main

import (
	"flag"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"crawler/pkg/crawler"
)

func main() {
	cfgPath := flag.String("c", "", "path to config file")
	dbPath := flag.String("d", "", "path to db file")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&crawler.Article{})

	crawler.ConfigCrawler(*cfgPath)
	crawler.Run(db)
}
