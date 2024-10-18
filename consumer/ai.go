package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type AI struct {
	Addr   string
	AppId  string
	Client *resty.Client
}

type Res struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

type chatResData struct {
	Content string `json:"content"`
}

func NewAI(addr, ak, appid string) *AI {
	client := resty.New()
	client.SetBaseURL(addr)
	client.SetHeader("Authorization", ak)
	return &AI{
		Addr:   addr,
		AppId:  appid,
		Client: client,
	}
}

func (a *AI) Send(chatid, text string) (string, error) {
	res, err := a.Client.R().
		SetBody(map[string]any{
			"message": text,
			"re_chat": true,
			"stream":  false,
		}).
		Post(fmt.Sprintf("/api/application/chat_message/%s", chatid))
	if err != nil {
		return "", err
	}

	resS := new(Res)
	err = json.Unmarshal(res.Body(), resS)
	if err != nil {
		return "", err
	}

	if resS.Code != 200 {
		return "", fmt.Errorf("code: %d, message: %s", resS.Code, resS.Message)
	}

	str := resS.Data.(map[string]any)["content"]
	return str.(string), nil
}

func (a *AI) GetChatID() (string, error) {
	res, err := a.Client.R().Get(fmt.Sprintf("/api/application/%s/chat/open", a.AppId))
	if err != nil {
		return "", err
	}

	resS := new(Res)
	err = json.Unmarshal(res.Body(), resS)
	if err != nil {
		return "", err
	}

	return resS.Data.(string), nil
}
