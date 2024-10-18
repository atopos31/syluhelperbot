package models

type Merchant struct {
	ID         int64   `json:"id"`
	Title      string  `json:"title"`
	Cover      string  `json:"thumb"`       // 封面
	Score      float32 `json:"score"`       // 评分
	Said       string  `json:"seller_said"` // 宣传语
	StartPrice float32 `json:"start_price"` // 起购价
	Status     int     `json:"open_status"` // 状态 1 关门 2 开门
}

func (m *Merchant) GetStatus() string {
	switch m.Status {
	case 1:
		return "关门"
	case 2:
		return "开门"
	default:
		return "未知"
	}
}
