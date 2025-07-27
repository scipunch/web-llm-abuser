package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	h := ResponsesHandler{
		userInput: make(chan string),
		wsMessage: make(chan WSMessage),
	}

	mux := http.NewServeMux()

	mux.Handle("GET /responses/ws", websocket.Handler(h.WebSocket))
	mux.HandleFunc("POST /responses", h.Post)

	host := os.Getenv("HOST")
	if host == "" {
		host = ":8080"
	}
	server := http.Server{
		Addr:    host,
		Handler: mux,
	}

	slog.Info("starting listening", "on", host)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type ResponsesHandler struct {
	userInput chan string
	wsMessage chan WSMessage
}

type PostRequest struct {
	Input string `json:"input"`
}

type PostResponse struct {
	Output string `json:"output"`
}

func (h ResponsesHandler) Post(w http.ResponseWriter, r *http.Request) {
	var reqBody PostRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if reqBody.Input == "" {
		http.Error(w, "input should not be empty", http.StatusUnprocessableEntity)
		return
	}
	slog.Info("responses request", "body", reqBody)

	h.userInput <- reqBody.Input
	wsResp := <-h.wsMessage

	var resp PostResponse
	switch wsResp.T {
	case "error":
		slog.Error("chromium extension failed", "with", wsResp.Data)
		http.Error(w, wsResp.Data, http.StatusInternalServerError)
		return
	case "model-output":
		slog.Info("got LLM response", "length", len(wsResp.Data))
		resp.Output = wsResp.Data
	default:
		log.Fatalf("unexpected chromium extension response %v", wsResp)
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

func (h ResponsesHandler) WebSocket(ws *websocket.Conn) {
	slog.Info("got new WebSocket connection")
	wsMessage := make(chan WSMessage)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			var msg WSMessage
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				slog.Error("could not receive", "with", err)
				cancel()
				break
			}
			slog.Info("received message from WebSocket")
			wsMessage <- msg
		}
	}()

	for {
		select {
		case input := <-h.userInput:
			err := websocket.JSON.Send(ws, WSMessage{T: "user-input", Data: input})
			if err != nil {
				slog.Error("could not send", "with", err)
				cancel()
			}
			slog.Info("input text sent via WebSocket")
		case wsMsg := <-wsMessage:
			slog.Info("received WebSocket message", "data", wsMsg)
			h.wsMessage <- wsMsg
		case <-ctx.Done():
			slog.Warn("socket with LLM closed")
			return
		}
	}
}

type WSMessage struct {
	T    string `json:"type"`
	Data string `json:"data"`
}
