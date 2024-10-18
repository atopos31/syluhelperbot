package consumer

import (
	"log"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestCum(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: "192.168.0.105:3001", Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()
}
