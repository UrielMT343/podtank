package ws

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	gameWS "github.com/UrielMT343/podtank/internal/serverWebSocket/game"
	"github.com/coder/websocket"
)

type Handler struct {
	RoomManager *gameWS.RoomManager
}

func NewHandler(rm *gameWS.RoomManager) *Handler {
	return &Handler{RoomManager: rm}
}

func (h *Handler) WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("Error while accepting ws connection", "error", err)
		return
	}
	defer conn.CloseNow()

	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*10)
	defer cancel()

	for {
		session := gameWS.PlayerSession{
			Conn:   conn,
			Ctx:    ctx,
			Cancel: cancel,
		}

		h.RoomManager.JoinOrCreate(&session)

		msgType, data, err := conn.Read(ctx)
		if err != nil {
			slog.Error("Connection closed", "error", err)
			break
		}

		slog.Info("Message received", "data", data)

		errWrite := conn.Write(ctx, msgType, data)
		if errWrite != nil {
			slog.Error("Error on writing message", "error", errWrite)
			break
		}
	}
}
