package main

import (
	"net/http"
	"strings"

	"github.com/UrielMT343/podtank/internal/server/signaling"
)

type Handlers struct {
	signalingHandler *signaling.Handler
}

func routes(h *Handlers, basePrefix string) http.Handler {
	apiMux := http.NewServeMux()

	signalingMux := http.NewServeMux()
	signalingMux.HandleFunc("POST /offer", h.signalingHandler.Offer)

	prefix := basePrefix
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	apiMux.Handle(prefix, http.StripPrefix(strings.TrimSuffix(prefix, "/"), signalingMux))

	return apiMux
}
