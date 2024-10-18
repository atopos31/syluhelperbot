package main

import (
	"bot/botcore"
	"bot/consumer"
	"bot/cron"
	"bot/listener"
	"bot/models"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	uin := os.Getenv("uin")
	wshost := os.Getenv("wshost")
	webapihost := os.Getenv("webapihost")
	AK := os.Getenv("AK")
	maxdbwebapiaddr := os.Getenv("maxdbwebapiaddr")
	maxdbak := os.Getenv("maxdbak")
	maxdbappid := os.Getenv("maxdbappid")
	var err error
	conn := new(websocket.Conn)
	var webuitoken string
	var trystart int
	for {
		conn, err = botcore.Connect(wshost)
		if err != nil {
			trystart++
			log.Println("Connect WS error:", err)
			time.Sleep(3 * time.Second)
			if strings.EqualFold(webuitoken, "") {
				webuitoken, err = botcore.GetWebUIToken(webapihost, AK)
				if err != nil {
					log.Println("Get webui token error:", err)
				}
			}
			if err := botcore.SetQuickLogin(webapihost, webuitoken, uin); err != nil {
				log.Println("Set quick login error:", err)
			}

			if trystart > 20 {
				restartNCcmd := exec.Command("docker", "restart", "napcat")
				_, err = restartNCcmd.CombinedOutput()
				log.Println("Restart nc failed:", err)
				return
			}
			continue
		}
		break
	}
	log.Println("Connect success")
	bot := botcore.NewBot(conn, uin)
	bot.SendPrivateMessage(2945294768, models.Message{
		Typ: "text",
		Data: models.Data{
			Text: "bot 已启动！",
		},
	})
	// 监听服务
	aimsgchan := make(chan models.Chanmsg)
	cmdmsgchan := make(chan models.Cmdmsg)
	Listener := listener.NewListener(aimsgchan, cmdmsgchan, bot)
	go Listener.Listen()
	// AI服务
	ai := consumer.NewAI(maxdbwebapiaddr, maxdbak, maxdbappid)
	sess := consumer.NewChatSession()
	// 吃什么服务
	merchhantMgr := cron.NewMerchantMgr()
	go cron.UpdateMerchantList(merchhantMgr)

	// 消费者
	Consumer := consumer.NewConsumer(aimsgchan, cmdmsgchan, ai, sess, merchhantMgr, bot)
	Consumer.Start()
}
