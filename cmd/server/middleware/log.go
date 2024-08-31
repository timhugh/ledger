package middleware

import (
	"log"
	"net/http"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s request_id=%s",
			r.Method,
			r.URL.Path,
			r.Context().Value("request_id"),
		)
	})
}
