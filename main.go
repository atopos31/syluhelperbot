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
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	uin := os.Getenv("uin")
	Masterqq := os.Getenv("masterqq")
	GroupIdstr := os.Getenv("groupid")
	models.GroupId, _ = strconv.ParseInt(GroupIdstr, 10, 64)
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
				trystart = 0
				restartNCcmd := exec.Command("docker", "restart", "napcat")
				_, err = restartNCcmd.CombinedOutput()
				if err != nil {
					log.Println("Restart nc failed:", err)
				}
			}
			continue
		}
		break
	}
	log.Println("Connect success")
	masterqqint, _ := strconv.ParseInt(Masterqq, 10, 64)
	bot := botcore.NewBot(conn, uin, masterqqint)
	bot.SendPrivateMessage(masterqqint, models.Message{
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

	//  races 服务
	raceMgr := cron.NewRaceMgr(bot)
	go raceMgr.Start()

	// 消费者
	Consumer := consumer.NewConsumer(aimsgchan, cmdmsgchan, ai, sess, merchhantMgr,raceMgr, bot)
	Consumer.Start()
}
