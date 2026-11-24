package signaling

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/UrielMT343/podtank/internal/server/response"
	"github.com/pion/webrtc/v4"
)

type signalReq struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Offer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Wrong HTTP Method")
		slog.Error("Wrong method passed", "method", r.Method)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var req signalReq
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if err := dec.Decode(&req); err != nil {
		m := response.FormatErrorMessage("Bad JSON", err.Error())
		response.Error(w, http.StatusBadRequest, m)
		return
	}
	if req.Type != "offer" || req.SDP == "" {
		response.Error(w, http.StatusBadRequest, "Expected type=offer with sdp")
		slog.Error("Wrong type with sdp", "type", req.Type)
		return
	}

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		slog.Error("Error while creating peer conn", "error", err.Error())
		return
	}
	defer pc.Close()

	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  req.SDP,
	}
	if err := pc.SetRemoteDescription(offer); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		slog.Error("Error while setting peer conn description", "error", err.Error())
		return
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error while creating answer")
		slog.Error("Error while creating answer", "error", err.Error())
		return
	}

	if err := pc.SetLocalDescription(answer); err != nil {
		response.Error(w, http.StatusInternalServerError, "Error while setting local description")
		slog.Error("Error while setting local decription", "error", err.Error())
		return
	}

	resp := signalReq{
		Type: answer.Type.String(),
		SDP:  answer.SDP,
	}

	response.Respond(w, http.StatusOK, resp)
}
