package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

var douyinPrefix = "https://live.douyin.com/webcast/room/web/enter/"

type DouyinStreamCrawler struct {
	RoomId     string `json:"roomId"`
	Name       string `json:"name"`
	NeedNotify bool
}

type DouyinRoomInfo struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

type DouyinData struct {
	Data []DouyinRoomInfo `json:"data"`
}

type DouyinResp struct {
	StatusCode int        `json:"status_code"`
	Data       DouyinData `json:"data"`
}

func (r *DouyinStreamCrawler) Crawl() error {
	params := url.Values{}
	params.Add("aid", "6383")
	params.Add("device_platform", "web")
	params.Add("enter_from", "web_live")
	params.Add("cookie_enabled", "true")
	params.Add("browser_language", "zh-CN")
	params.Add("browser_platform", "Win32")
	params.Add("browser_name", "Chrome")
	params.Add("browser_version", "109.0.0.0")
	params.Add("web_rid", r.RoomId)

	douyinUrl := fmt.Sprintf("%s?%s", douyinPrefix, params.Encode())

	response, err := http.Get(douyinUrl)
	if err != nil {
		logrus.Errorf("get room [%s] status with url [%s] failed: %s", r.RoomId, douyinUrl, err.Error())
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

	var resp DouyinResp
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		logrus.Errorf("unmarshal response body failed: %s", err.Error())
		return err
	}

	if len(resp.Data.Data) == 0 {
		err = fmt.Errorf("response data len < 1: %v", resp)
		logrus.Error(err.Error())
		return err
	}

	if resp.Data.Data[0].Status == 2 {
		if r.NeedNotify {
			logrus.Infof("room [%s] start stream", r.Name)
			textChan <- fmt.Sprintf("[%s]@Douyin开播: %s", r.Name, resp.Data.Data[0].Title)
			r.NeedNotify = false
		}
	} else {
		logrus.Infof("room [%s] has not stream", r.Name)
		r.NeedNotify = true
	}
	return nil
}
