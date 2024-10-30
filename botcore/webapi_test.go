package botcore

import "testing"

const webapihost = "192.168.0.105:6099"
const AK = "tgglqkaesto"

func TestWebApi(t *testing.T) {
	cer, err := GetWebUIToken(webapihost, AK)
	if err != nil {
		t.Error(err)
	}
	t.Log(cer)
	t.Log(CheckLoginStatus(webapihost, cer))
	if err := SetQuickLogin(webapihost, cer, "3808139675"); err != nil {
		t.Log(err)
	}
}
