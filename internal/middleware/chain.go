package middleware

import (
	"net/http"
	"slices"
)

func CreateChain(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range slices.Backward(middlewares) {
		next = m(next)
	}

	return next
}
