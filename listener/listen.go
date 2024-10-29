package listener

import (
	"bot/botcore"
	"bot/models"
	"log"
	"strconv"
	"strings"
)

type Listener struct {
	Aimsgchan  chan models.Chanmsg
	Cmdmsgchan chan models.Cmdmsg
	Bot        *botcore.Bot
}

func NewListener(aimsgchan chan models.Chanmsg, cmdmsgchan chan models.Cmdmsg, bot *botcore.Bot) *Listener {
	return &Listener{Aimsgchan: aimsgchan, Cmdmsgchan: cmdmsgchan, Bot: bot}
}

func (l *Listener) Listen() {
	for {
		message, err := l.Bot.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}

		log.Println("message:", message)

		go l.Handler(message)
	}
}

func (l *Listener) Handler(msg *models.MessageData) {
	if msg.MessageType != "group" || msg.GroupID != models.GroupId {
		return
	}
	var atqq string
	var text string
	for data := range msg.Message {
		if msg.Message[data].Typ == "at" {
			atqq = msg.Message[data].Data.QQ
		} else if msg.Message[data].Typ == "text" {
			text = msg.Message[data].Data.Text
		}
	}

	qq := strconv.FormatInt(msg.UserID, 10)
	log.Println("QQ:", msg.UserID, "Text:", text)

	if strings.Contains(text, "/") {
		cmdmsg := models.Cmdmsg{
			Cmd: text,
		}
		l.Cmdmsgchan <- cmdmsg
	} else if strings.EqualFold(qq, atqq) {
		chanmsg := models.Chanmsg{
			QQ:   qq,
			Text: text,
		}
		l.Aimsgchan <- chanmsg
	}
}
