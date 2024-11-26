package limitedMiddlewarepackage

import (
	"net/http"
)

var sem = make(chan struct{}, 10)

func LimitedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sem <- struct{}{}
		defer func() { <-sem }()
		next.ServeHTTP(w, r)
	})
}
