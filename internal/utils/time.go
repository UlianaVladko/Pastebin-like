package utils

import (
	"strconv"
	"time"
)

func TimeAgo(t time.Time) string {
	d := time.Since(t)

	minutes := int(d.Minutes())
	hours := int(d.Hours())
	days := int(d.Hours() / 24)
	weeks := days / 7
	months := days / 30
	years := days / 365

	if d < time.Minute {
		return "just now"
	} else if d < time.Hour {
		if minutes == 1 {
			return "1 minute ago"
		}
		return strconv.Itoa(minutes) + " minutes ago"
	} else if d < 24*time.Hour {
		if hours == 1 {
			return "1 hour ago"
		}
		return strconv.Itoa(hours) + " hours ago"
	} else if days < 7 {
		if days == 1 {
			return "1 day ago"
		}
		return strconv.Itoa(days) + " days ago"
	} else if days < 30 {
		if weeks == 1 {
			return "1 week ago"
		}
		return strconv.Itoa(weeks) + " weeks ago"
	} else if days < 365 {
		if months == 1 {
			return "1 month ago"
		}
		return strconv.Itoa(months) + " months ago"
	} else {
		if years == 1 {
			return "1 year ago"
		}
		return strconv.Itoa(years) + " years ago"
	}
}
