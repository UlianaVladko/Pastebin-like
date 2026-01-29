package handlers

import (
	"strconv"
	"time"
)

func computeExpiration(exp string) *time.Time {
	now := time.Now().UTC()
	var t time.Time

	switch exp {
	case "10m":
		t = now.Add(10 * time.Minute)
	case "1h":
		t = now.Add(time.Hour)
	case "1d":
		t = now.Add(24 * time.Hour)
	default:
		return nil
	}
	return &t
}

func timeAgo(t time.Time) string {
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
