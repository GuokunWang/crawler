package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

var prefix = "https://api.live.bilibili.com/room/v1/Room/get_info?room_id="

type BiliStreamCrawler struct {
	RoomId string `json:"roomId"`
	Name   string `json:"name"`
}

type RoomInfo struct {
	Uid        int    `json:"uid"`
	LiveStatus int    `json:"live_status"`
	Title      string `json:"title"`
}

type Resp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Msg     string   `json:"msg"`
	Data    RoomInfo `json:"data"`
}

func (r *BiliStreamCrawler) Crawl() error {
	url := fmt.Sprintf("%s%s", prefix, r.RoomId)
	response, err := http.Get(url)
	if err != nil {
		logrus.Errorf("get room [%s] status with url [%s] failed: %s", r.RoomId, url, err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logrus.Errorf("error: status code %d", response.StatusCode)
		return fmt.Errorf("error: status code %d", response.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("read response body failed: %s", err.Error())
		return err
	}

	var resp Resp
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		logrus.Errorf("unmarshal response body failed: %s", err.Error())
		return err
	}

	if resp.Data.LiveStatus == 1 {
		logrus.Infof("room [%s] start stream", r.Name)
		textChan <- fmt.Sprintf("[%s]@BiliBili开播: %s", r.Name, resp.Data.Title)
	} else {
		logrus.Infof("room [%s] has not stream", r.Name)
	}
	return nil
}
