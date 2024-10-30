package cron

import (
	"bot/models"
	"errors"
	"sync"

	"golang.org/x/exp/rand"
)

const MerchantOssDomain = "https://shizaixiaoyuan.oss-cn-huhehaote.aliyuncs.com"

type MerchantMgr struct {
	mu           sync.Mutex
	MerchantList []models.Merchant
}

func NewMerchantMgr() *MerchantMgr {
	return &MerchantMgr{
		MerchantList: make([]models.Merchant, 0, 200),
	}
}

// 随机获取一个
func (m *MerchantMgr) GetRandomMerchant() (*models.Merchant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	num := len(m.MerchantList)
	if num == 0 {
		return nil, errors.New("empty merchant list")
	}
	return &m.MerchantList[rand.Intn(num)], nil
}

func (m *MerchantMgr) Update(merchant []models.Merchant) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MerchantList = merchant
}

func (m *MerchantMgr) GetMerchantList() []models.Merchant {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.MerchantList
}
