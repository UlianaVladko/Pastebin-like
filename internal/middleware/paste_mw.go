package middleware

import (
	"context"
	"net/http"
	"pastebin/internal/models"
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
