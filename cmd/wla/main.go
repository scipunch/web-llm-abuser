package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	userMessages := make(chan Message)
	llmMessages := make(chan Message)

	mux := http.NewServeMux()

	mux.Handle("GET /websocket", websocket.Handler(LLMExchange(userMessages, llmMessages)))
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var msg Message
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userMessages <- msg
		msg = <-llmMessages

		respJSON, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func LLMExchange(req <-chan Message, resp chan Message) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		llmAnswer := make(chan Message)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			for {
				var msg Message
				if err := websocket.JSON.Receive(ws, &msg); err != nil {
					slog.Error("could not receive", "with", err)
					cancel()
					break
				}
				llmAnswer <- msg
			}
		}()

		for {
			select {
			case userMsg := <-req:
				err := websocket.JSON.Send(ws, userMsg)
				if err != nil {
					slog.Error("could not send", "with", err)
					cancel()
				}
			case llmMsg := <-llmAnswer:
				slog.Info("received LLM message", "length", len(llmMsg.Message))
				resp <- llmMsg
			case <-ctx.Done():
				slog.Warn("socket with LLM closed")
				return
			}
		}
	}
}

type Message struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}
