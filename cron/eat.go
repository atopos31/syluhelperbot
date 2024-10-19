package cron

import (
	"bot/models"
	"encoding/json"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

type res struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Items []models.Merchant `json:"items"`
	} `json:"data"`
}

func UpdateMerchantList(mgr *MerchantMgr) {
	client := resty.New()
	// 定时任务
	ticker := time.NewTicker(time.Second * 2)
	merchants, err := getMerchantList(client)
	if err != nil {
		log.Println("获取商家列表失败", err)
	} else {
		mgr.Update(merchants)
	}
	for range ticker.C {
		merchants, err := getMerchantList(client)
		if err != nil {
			log.Println("获取商家列表失败", err)
			continue
		}
		mgr.Update(merchants)
	}
}

func getMerchantList(client *resty.Client) ([]models.Merchant, error) {
	resp, err := client.R().
		SetQueryParams(map[string]string{"lat": "41.730958", "lng": "123.502544", "skip": "0"}).
		SetHeaders(map[string]string{"x-appid": "Lk4X5qCT0d1oX1DC", "x-model": "PC", "x-platform": "windows"}).
		Get("https://a.shizaixiaoyuan.cn/api/merchant/merchantList")
	if err != nil {
		return nil, err
	}
	resstruct := new(res)
	if err = json.Unmarshal(resp.Body(), resstruct); err != nil {
		return nil, err
	}
	return resstruct.Data.Items, nil
}
