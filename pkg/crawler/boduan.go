package crawler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type BoduanCrawler struct {
	Url       string `json:"url"`
	Cookie    string `json:"cookie"`
	Token     string `json:"token"`
	FakeId    string `json:"fakeId"`
}

type RespItem struct {
	Title string `json:"title"`
	Link string `json:"link"`
}

type Resp struct {
	AppMsgCnt int `json:"app_msg_cnt"`
	AppMsgList []RespItem`json:"app_msg_list"`
}


func (b *BoduanCrawler) Crawl() error {
	headers := http.Header{
		"Cookie":     []string{b.Cookie},
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36"},
	}

	params := map[string]string{
		"token":  b.Token,
		"lang":   "zh_CN",
		"f":      "json",
		"ajax":   "1",
		"action": "list_ex",
		"begin":  "0",
		"count":  "5",
		"query":  "",
		"fakeid": b.FakeId,
		"type":   "9",
	}

	req, err := http.NewRequest("GET", b.Url, nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Decode the JSON response into a map.
	content := Resp{}
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		panic(err)
	}

	// Extracting the title and corresponding url of each page article
	for _, appMsg := range content.AppMsgList {
		article := Article{
			ID:    uuid.New().String(),
			Title: appMsg.Title,
			URL:   appMsg.Link,
		}
		articleChan <- article
	}
	return nil
}