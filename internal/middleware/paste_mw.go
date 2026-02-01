package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"pastebin/internal/models"
	"strings"
	"time"
)

func PasteMiddleware(
	loader func(*http.Request) (*models.Paste, error),
	render404 func(http.ResponseWriter),
	next http.HandlerFunc,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		paste, err := loader(r)
		if err != nil || paste == nil {
			render404(w)
			return
		}

		if paste.ExpiresAt != nil && time.Now().UTC().After(*paste.ExpiresAt) {
			render404(w)
			return
		}

		ctx := context.WithValue(r.Context(), PasteKey, paste)
		next(w, r.WithContext(ctx))
	}
}

func LoadPasteByShort(db *sql.DB) func(r *http.Request) (*models.Paste, error) {
	return func(r *http.Request) (*models.Paste, error) {
		short := strings.TrimPrefix(r.URL.Path, "/p/")
		if short == "" {
			return nil, sql.ErrNoRows
		}
		return models.GetPasteByShortLink(db, short)
	}
}

// func LoadPasteByShortPrefix(db *sql.DB, prefix string) func(r *http.Request) (*models.Paste, error) {
// 	return func(r *http.Request) (*models.Paste, error) {
// 		short := strings.TrimPrefix(r.URL.Path, prefix)
// 		short = strings.TrimPrefix(short, "/")
// 		if short == "" {
// 			return nil, sql.ErrNoRows
// 		}
// 		return models.GetPasteByShortLink(db, short)
// 	}
// }