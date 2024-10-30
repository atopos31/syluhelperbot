package util

import "testing"

func TestTime(t *testing.T) {
	ok, err := IsToday("2024-10-31 00:00:00")
	if err != nil {
		t.Error(err)
	}
	t.Log(ok)
}
