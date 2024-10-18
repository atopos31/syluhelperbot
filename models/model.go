package models


var GroupId int64

type MessageData struct {
	SelfID        int64      `json:"self_id"`
	UserID        int64      `json:"user_id"`
	Time          int64      `json:"time"`
	MessageID     int64      `json:"message_id"`
	MessageSeq    int64      `json:"message_seq"`
	RealID        int64      `json:"real_id"`
	MessageType   string     `json:"message_type"`
	Sender        SenderData `json:"sender"`
	RawMessage    string     `json:"raw_message"`
	Font          int        `json:"font"`
	SubType       string     `json:"sub_type"`
	Message       []Message  `json:"message"`
	MessageFormat string     `json:"message_format"`
	PostType      string     `json:"post_type"`
	GroupID       int64      `json:"group_id"`
}

type SenderData struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card"`
	Role     string `json:"role"`
}

type Message struct {
	Typ  string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	QQ   string `json:"qq,omitempty"`   // @ 某人
	Text string `json:"text,omitempty"` // 纯文本
	File string `json:"file,omitempty"` // 图片
}

type API struct {
	Action string `json:"action"`
	Params any    `json:"params"`
}

type ResGroup struct {
	GroupID int64     `json:"group_id"`
	Message []Message `json:"message"`
}

type ResPrivate struct {
	UserID  int64     `json:"user_id"`
	Message []Message `json:"message"`
}
