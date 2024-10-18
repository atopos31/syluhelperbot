package botcore

import (
	"bot/models"
	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"
)

type Bot struct {
	wsconn *websocket.Conn
	Mineqq string
}

func NewBot(conn *websocket.Conn,qq string) *Bot {
	return &Bot{
		wsconn: conn,
		Mineqq: qq,
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
