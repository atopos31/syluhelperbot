package cron

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
)

const TinyURLAPI = "https://t.hackerxiao.online"

type res struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"message"`
}

func GenTinyURL(URL string) (string,error) {
	client := resty.New()
	resp, err := client.R().
			SetBody(map[string]any{"url": URL,"time": 7}).
		Post(TinyURLAPI+"/create")
	if err != nil {
		return "", err
	}

	resstruct := new(res)
	if err = json.Unmarshal(resp.Body(), resstruct); err != nil {
		return "", err
	}
	log.Printf("tinyurl: %s",resstruct.Data)
	return resstruct.Data, nil
}