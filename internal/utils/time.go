package utils

import (
	"strconv"
	"time"
)

func TimeAgo(t time.Time) string {
	d := time.Since(t)
	
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return strconv.Itoa(int(d.Minutes())) + " minutes ago"
	case d < 24*time.Hour:
		return strconv.Itoa(int(d.Hours())) + " hours ago"
	default:
		return t.Format("02 Jan 2006")
	}
}