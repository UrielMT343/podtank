package main

import (
	"net/http"
	"strings"

	"github.com/UrielMT343/podtank/internal/serverWebSocket/ws"
)

type Handlers struct {
	wsHandler *ws.Handler
}

func routes(h *Handlers, basePrefix string) http.Handler {
	apiMux := http.NewServeMux()

	gameMux := http.NewServeMux()
	gameMux.HandleFunc("GET /game/ws", h.wsHandler.WsHandler)

	prefix := basePrefix
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	apiMux.Handle(prefix, http.StripPrefix(strings.TrimSuffix(prefix, "/"), gameMux))

	return apiMux
}
