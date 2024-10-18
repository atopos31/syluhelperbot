package consumer

import "testing"

func TestAi(t *testing.T) {
	ai := NewAI("http://192.168.0.105:8080", "application-f505694400b71c390e1108750ca5ce2f", "879c1fbe-7fce-11ef-a93d-0242ac110003")
	chatid, _ := ai.GetChatID()
	text, _ := ai.Send(chatid, "你好")
	t.Log(text)
}
