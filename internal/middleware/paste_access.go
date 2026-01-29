package middleware

import (
	"net/http"
	"pastebin/internal/models"
)

func PasteAccessMiddleware(
	render404 func(http.ResponseWriter),
	next http.HandlerFunc,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		paste, ok := r.Context().Value(PasteKey).(*models.Paste)
		if !ok || paste == nil {
			render404(w)
			return
		}

		if !paste.IsPrivate {
			next(w, r)
			return
		}

		token := r.URL.Query().Get("token")

		if token == "" || token != paste.EditToken {
			render404(w)
			return
		}

		next(w, r)
	}
}
