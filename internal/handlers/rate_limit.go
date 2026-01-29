package handlers

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	count     int
	lastReset time.Time
}

var (
	visitors   = make(map[string]*visitor)
	visitorsMu sync.Mutex
)

const (
	maxRequests = 5
	window      = time.Minute
)

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Cannot detect IP", http.StatusInternalServerError)
			return
		}

		visitorsMu.Lock()
		v, ok := visitors[ip]
		if !ok || time.Since(v.lastReset) > window {
			visitors[ip] = &visitor{
				count:     1,
				lastReset: time.Now(),
			}
			visitorsMu.Unlock()
			next(w, r)
			return
		}

		if v.count >= maxRequests {
			visitorsMu.Unlock()
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		v.count++
		visitorsMu.Unlock()

		next(w, r)
	}
}
