package cron

import "testing"

func TestEat(t *testing.T) {
	mgr := NewMerchantMgr()
	UpdateMerchantList(mgr)
}
