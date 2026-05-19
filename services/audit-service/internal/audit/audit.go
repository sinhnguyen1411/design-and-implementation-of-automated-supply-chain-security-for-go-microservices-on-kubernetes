package audit

import (
	"github.com/google/uuid"
	"time"
)

type Level string

const (
	LevelInfo Level = "info"
	LevelWarn Level = "warn"
	LevelCrit Level = "critical"
)

type Event struct {
	ID     string `json:"id"`
	Actor  string `json:"actor"`
	Action string `json:"action"`
	Level  Level  `json:"level"`
	TS     int64  `json:"ts"`
}

func NewEvent(actor, action string, level Level) Event {
	return Event{ID: uuid.NewString(), Actor: actor, Action: action, Level: level, TS: time.Now().UnixMilli()}
}
