package services

import (
	"database/sql"
	"log"
	"pastebin/internal/models"
	"time"
)

func StartExpiredPastesCleanup(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			if err := models.DeleteExpiredPastes(db); err != nil {
				log.Println("cleanup expired pastes error:", err)
			}
		}
	}()
}
