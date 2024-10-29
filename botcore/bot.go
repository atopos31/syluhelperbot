package botcore

import (
	"bot/models"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Bot struct {
	wsconn *websocket.Conn
	Mineqq string
	Masterqq int64
}

func NewBot(conn *websocket.Conn,mineqq string,masterqq int64) *Bot {
	return &Bot{
		wsconn: conn,
		Mineqq: mineqq,
		Masterqq: masterqq,
	}
}

func Connect(host string) (*websocket.Conn, error) {
	url := url.URL{Scheme: "ws", Host: host, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	return conn, err
}

func (bot *Bot) ReadMessage() (*models.MessageData, error) {
	_, rawmsg, err := bot.wsconn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var msg models.MessageData
	log.Println(string(rawmsg))
	return &msg, json.Unmarshal(rawmsg, &msg) // json转换可能会失败
}

func (bot *Bot) SendPrivateMessage(qq int64, msgs ...models.Message) error {
	msg := models.ResPrivate{
		UserID:  qq,
		Message: msgs,
	}
	return bot.wsconn.WriteJSON(&models.API{ 
		Action: "send_private_msg",
		Params: msg,
	})
}

func (bot *Bot) SendGroupMessage(groupid int64, msgs ...models.Message) error {
	msg := models.ResGroup{
		GroupID: groupid,
		Message: msgs,
	}
	return bot.wsconn.WriteJSON(&models.API{
		Action: "send_group_msg",
		Params: msg,
	})
}

func (bot *Bot) SendErrorMessage(err error) error {
	return bot.SendPrivateMessage(bot.Masterqq, models.Message{
		Typ: "text",
		Data: models.Data{
			Text: fmt.Sprintf("Error :%v", err),
		},
	})
}
