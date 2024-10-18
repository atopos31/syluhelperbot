package botcore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type res struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Credential string `json:"Credential"`
		IsLogin    bool   `json:"isLogin"`
	} `json:"data"`
}

func SetQuickLogin(host, credential, uin string) error {
	client := resty.New()
	url := fmt.Sprintf("http://%s%s", host, "/api/QQLogin/SetQuickLogin")
	resp, err := client.SetTimeout(3*time.Second).R().SetAuthToken(credential).SetBody(map[string]string{"uin": uin}).Post(url)
	if err != nil {
		return err
	}
	resstruct := new(res)
	if err = json.Unmarshal(resp.Body(), resstruct); err != nil {
		return err
	}
	if resstruct.Code != 0 {
		return fmt.Errorf("set quick login failed: %s", resstruct.Message)
	}

	return nil
}

func CheckLoginStatus(host, credential string) (bool, error) {
	client := resty.New()
	url := fmt.Sprintf("http://%s%s", host, "/api/QQLogin/CheckLoginStatus")
	resp, err := client.R().SetAuthToken(credential).Post(url)
	if err != nil {
		return false, err
	}
	resstruct := new(res)
	if err = json.Unmarshal(resp.Body(), resstruct); err != nil {
		return false, err
	}

	return resstruct.Data.IsLogin, nil
}

func GetWebUIToken(host, ak string) (string, error) {
	client := resty.New()
	url := fmt.Sprintf("http://%s%s", host, "/api/auth/login")
	resp, err := client.SetTimeout(3*time.Second).R().SetBody(map[string]string{"token": ak}).Post(url)
	if err != nil {
		return "", err
	}

	resstruct := new(res)
	if err = json.Unmarshal(resp.Body(), resstruct); err != nil {
		return "", err
	}
	return resstruct.Data.Credential, nil
}
