package health

import (
	"regexp"
)

var rule = regexp.MustCompile("^[a-z0-9-]+$")

func Status() string {
	if rule.MatchString("risk-service") {
		return "ok:rules"
	}
	return "invalid"
}
