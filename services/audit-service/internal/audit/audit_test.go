package audit

import "testing"

func TestNewEvent(t *testing.T) {
	e := NewEvent("user-1", "login", LevelInfo)
	if e.ID == "" { t.Fatal("empty id") }
	if e.Actor != "user-1" { t.Fatal("wrong actor") }
	if e.TS == 0 { t.Fatal("zero timestamp") }
}
