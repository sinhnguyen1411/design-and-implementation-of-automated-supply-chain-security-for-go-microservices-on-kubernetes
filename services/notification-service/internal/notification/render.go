package notification

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Request struct {
	EventType string            `json:"event_type"`
	Severity  string            `json:"severity"`
	Actor     string            `json:"actor"`
	Target    string            `json:"target"`
	Metadata  map[string]string `json:"metadata"`
}

type Message struct {
	MessageID string `json:"message_id"`
	Channel   string `json:"channel"`
	Priority  string `json:"priority"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Valid     bool   `json:"valid"`
	Reason    string `json:"reason,omitempty"`
}

func Render(req Request) Message {
	eventType := strings.ToLower(strings.TrimSpace(req.EventType))
	severity := strings.ToLower(strings.TrimSpace(req.Severity))
	out := Message{MessageID: uuid.NewString()}

	if eventType == "" || severity == "" {
		out.Valid = false
		out.Reason = "invalid_payload"
		return out
	}

	channel := "email"
	priority := "p3"
	switch severity {
	case "critical":
		channel = "pagerduty"
		priority = "p1"
	case "high":
		channel = "slack"
		priority = "p2"
	case "medium", "low", "info":
	default:
		out.Valid = false
		out.Reason = "invalid_severity"
		return out
	}

	title := fmt.Sprintf("[%s] %s", strings.ToUpper(severity), strings.ToUpper(eventType))
	body := fmt.Sprintf("actor=%s target=%s", req.Actor, req.Target)
	if runID, ok := req.Metadata["run_id"]; ok && strings.TrimSpace(runID) != "" {
		body += fmt.Sprintf(" run_id=%s", runID)
	}

	out.Channel = channel
	out.Priority = priority
	out.Title = title
	out.Body = body
	out.Valid = true
	return out
}

