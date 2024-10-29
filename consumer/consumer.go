package consumer

import (
	"bot/botcore"
	"bot/cron"
	"bot/models"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Consumer struct {
	ToAIMsg    chan models.Chanmsg
	Cmdmsg     chan models.Cmdmsg
	Ai         *AI
	Sess       *chatSession
	MerchatMgr *cron.MerchantMgr
	RaceMgr    *cron.RaceMgr
	Bot        *botcore.Bot
}

func NewConsumer(aimsg chan models.Chanmsg, cmdmsg chan models.Cmdmsg, ai *AI, sess *chatSession, mgr *cron.MerchantMgr,raceMgr *cron.RaceMgr, bot *botcore.Bot) *Consumer {
	return &Consumer{ToAIMsg: aimsg, Cmdmsg: cmdmsg, Ai: ai, Sess: sess, MerchatMgr: mgr,RaceMgr: raceMgr, Bot: bot}
}

func (c *Consumer) Start() {
	for {
		select {
		case msg := <-c.ToAIMsg:
			log.Printf("msg %v", msg)
			if err := c.SettleAI(msg); err != nil {
				log.Printf("settle ai failed %v", err)
				c.Bot.SendErrorMessage(err)
			}
		case cmd := <-c.Cmdmsg:
			log.Printf("cmd %v", cmd)
			if err := c.SettleCmd(cmd); err != nil {
				log.Printf("settle cmd failed %v", err)
				c.Bot.SendErrorMessage(err)
			}
		}
	}
}

func (c *Consumer) SettleAI(msg models.Chanmsg) error {
	var err error
	chatid, ok := c.Sess.Get(msg.QQ)
	if !ok {
		chatid, err = c.Ai.GetChatID()
		if err != nil {
			return err
		}
		c.Sess.Set(msg.QQ, chatid)
	}

	res, err := c.Ai.Send(chatid, msg.Text)
	if err != nil {
		return err
	}
	return c.Bot.SendGroupMessage(models.GroupId, models.Message{
		Typ: "at",
		Data: models.Data{
			QQ: msg.QQ,
		},
	}, models.Message{
		Typ: "text",
		Data: models.Data{
			Text: fmt.Sprintf(" %s", res),
		},
	})
}

func (c *Consumer) SettleCmd(cmd models.Cmdmsg) error {
	switch {
		case strings.Contains(cmd.Cmd, "help"):
			return c.cmdHelp()
		case strings.Contains(cmd.Cmd, "今天吃什么"):
			return c.cmdEat()
		case strings.Contains(cmd.Cmd, "近期比赛"):
			return c.cmdRaces()
		default:
			return errors.New("cmd not found")
	}
}

func (c *Consumer) cmdHelp() error {
	return c.Bot.SendGroupMessage(models.GroupId,
		models.Message{
		Typ: "text",
		Data: models.Data{
			Text: "/今天吃什么：随机推荐食堂商家\n",
		},
	},models.Message{
		Typ: "text",
		Data: models.Data{
			Text: "/近期比赛：查询近期竞赛",
		},
	})
}

func (c *Consumer) cmdEat() error {
	merchant, err := c.MerchatMgr.GetRandomMerchant()
	if err != nil {
		return err
	}
	if strings.EqualFold(merchant.Said, "") {
		merchant.Said = "小助手推荐！"
	}
	return c.Bot.SendGroupMessage(models.GroupId,
		models.Message{
			Typ: "image",
			Data: models.Data{
				File: fmt.Sprintf("%s/%s", cron.MerchantOssDomain, merchant.Cover),
			},
		}, models.Message{
			Typ: "text",
			Data: models.Data{
				Text: fmt.Sprintf("今天吃%s吧，%s\n", merchant.Title, merchant.Said),
			},
		}, models.Message{
			Typ: "text",
			Data: models.Data{
				Text: fmt.Sprintf("评分：%f 状态：%s", merchant.Score, merchant.GetStatus()),
			},
		})
}

func (c *Consumer) cmdRaces() error {
	c.RaceMgr.Mu.Lock()
	defer c.RaceMgr.Mu.Unlock()
	msgs := make([]models.Message,0)
	for _,race := range c.RaceMgr.Races {
		msgs = append(msgs,models.Message{
			Typ: "text",
			Data: models.Data{
				Text: fmt.Sprintf("%s\n%s\n",race.Name,race.URL),
			},
		})
	}
	msgs[len(msgs)-1].Data.Text = strings.TrimSuffix(msgs[len(msgs)-1].Data.Text,"\n")
	return c.Bot.SendGroupMessage(models.GroupId,msgs...)
}
