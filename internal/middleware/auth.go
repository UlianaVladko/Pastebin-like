package middleware

import "net/http"

func AuthMiddleware(render401 func(http.ResponseWriter), next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			render401(w)
			return
		}
		next(w, r)
	}
}

func Render401(w http.ResponseWriter)  {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("401 Unauthorized: Please log in"))
}
