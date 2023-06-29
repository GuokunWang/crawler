package crawler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type DingTalkClient struct {
	WebHookUrl string `json:"webHookUrl"`
	Secret     string `json:"secret"`
}

func (d *DingTalkClient) PushArticles(articles []Article) {
	links := []map[string]interface{}{}
	for _, article := range articles {
		link := map[string]interface{}{}
		link["title"] = article.Title
		link["messageURL"] = article.URL
		link["picURL"] = ""
		links = append(links, link)

	}
	message := map[string]interface{}{
		"msgtype": "feedCard",
		"feedCard": map[string]interface{}{
			"links": links,
		},
	}

	jsonMessage, _ := json.Marshal(message)
	url := d.ModifyURL()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}

func (d *DingTalkClient) ModifyURL() string {
	secret := []byte(d.Secret)
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	signStr := fmt.Sprintf("%s\n%s", timestamp, secret)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(sign)
	url := fmt.Sprintf("%s&timestamp=%s&sign=%s", d.WebHookUrl, timestamp, sign)
	return url
}
