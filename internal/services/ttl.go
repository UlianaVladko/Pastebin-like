package services

import "time"

func ComputeExpiration(exp string) *time.Time {
	now := time.Now().UTC()

	switch exp {
	case "10m":
		t := now.Add(10 * time.Minute)
		return &t
	case "1h":
		t := now.Add(time.Hour)
		return &t
	case "1d":
		t := now.Add(24 * time.Hour)
		return &t
	default:
		return nil
	}
}