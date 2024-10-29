package cron

import (
	"testing"
)

func TestRace(t *testing.T) {
	
	NewRaceMgr(nil).Start()
}