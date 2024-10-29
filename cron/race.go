package cron

import (
	"bot/botcore"
	"bot/models"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const DefaultCoverPhoto = "https://cxcy.upln.cn/img/bg3.9cc9e2c5.jpg"

type resRace struct {
	Code   int    `json:"code"`
	Msg    string `json:"messsage"`
	Result struct {
		Races []Race `json:"records"`
	} `json:"result"`
}

type RaceMgr struct {
	Mu           sync.Mutex
	Bot   *botcore.Bot
	Races map[string]Race
}

type Race struct {
	Name       string `json:"name"`
	CoverPhoto string `json:"coverPhoto"`
	ID         string `json:"id"`
	EndTime    string `json:"endTime"`
	URL        string
}

func NewRaceMgr(bot *botcore.Bot) *RaceMgr {
	return &RaceMgr{Bot: bot,Races: map[string]Race{}}
}

func (r *RaceMgr) Start() {
	client := resty.New()
	races, err := r.getNewRaces(client)
	if err != nil {
		log.Println(err)
		return
	}

	for _, race := range races {
		race.URL, err = GenTinyURL("https://cxcy.upln.cn/match/details?comId=" + race.ID)
		if err != nil {
			r.Bot.SendErrorMessage(err)
		}
		log.Println(race.URL)
		r.Mu.Lock()
		r.Races[race.ID] = race
		r.Mu.Unlock()
	}
	log.Println("race success:")
	log.Println(r.Races)
	t := time.NewTicker(time.Hour)
	for range t.C {
		races, err := r.getNewRaces(client)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, race := range races {
			race.URL, err = GenTinyURL("https://cxcy.upln.cn/match/details?comId=" + race.ID)
			if err != nil {
				r.Bot.SendErrorMessage(err)
			}
			r.Mu.Lock()
			if _, ok := r.Races[race.ID]; !ok {
				if err:=r.SendNewRace(race); err!=nil {
					r.Bot.SendErrorMessage(err)
				}
			}
			r.Races[race.ID] = race
			r.Mu.Unlock()
		}
	}
}

func (r *RaceMgr) SendNewRace(race Race) error {
	if strings.EqualFold(race.CoverPhoto, "") {
		race.CoverPhoto = DefaultCoverPhoto
	}
	msgs := []models.Message{{
		Typ: "text",
		Data: models.Data{
			Text: "有新的竞赛发布啦！",
		},
	}, {
		Typ: "image",
		Data: models.Data{
			File: race.CoverPhoto,
		},
	}, {
		Typ: "text",
		Data: models.Data{
			Text: "竞赛名称：" + race.Name + "\n",
		},
	}, {
		Typ: "text",
		Data: models.Data{
			Text: "报名截止时间：" + race.EndTime + "\n",
		},
	}, {
		Typ: "text",
		Data: models.Data{
			Text: "详情链接：" + race.URL + "\n",
		},
	}, {
		Typ: "text",
		Data: models.Data{
			Text: "小助手期待你的表现哦！",
		},
	}}
	if err := r.Bot.SendPrivateMessage(r.Bot.Masterqq, msgs...); err != nil {
		return err
	}
	return r.Bot.SendGroupMessage(models.GroupId, msgs...)
}

func (r *RaceMgr) getNewRaces(client *resty.Client) ([]Race, error) {
	// 请求接口获取 race 数据
	resp, err := client.R().
		SetQueryParams(map[string]string{"year": "2024", "code": "1", "pageNo": "1", "pageSize": "10"}).
		Get("https://cxcy.upln.cn/provincial/match/competition/queryOngoing")
	if err != nil {
		return nil, err
	}
	ress := new(resRace)
	if err := json.Unmarshal(resp.Body(), &ress); err != nil {
		return nil, err
	}
	return ress.Result.Races, nil
}

