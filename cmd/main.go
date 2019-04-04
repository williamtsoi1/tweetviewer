package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudevents/sdk-go"
	"github.com/mchmarny/tevents/pkg/handlers"
	"github.com/mchmarny/tevents/pkg/utils"
	"golang.org/x/net/websocket"
)

func main() {

	port, err := strconv.Atoi(utils.MustGetEnv("PORT", "8080"))
	if err != nil {
		log.Fatalf("failed to parse port, %s", err.Error())
	}

	// Configs
	handlers.InitHandlers()

	// Handler Mux
	mux := http.NewServeMux()

	// Static
	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	// Ingres API Handler
	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithMethod("POST"),
		cloudevents.WithPath("/"),
		cloudevents.WithPort(port),
	)
	if err != nil {
		log.Fatalf("failed to create cloudevents transport, %s", err.Error())
	}

	// wire handler for CE
	t.SetReceiver(&handlers.TwitterReceiver{})

	// WS Handler
	mux.Handle("/ws", websocket.Handler(handlers.WSHandler))
	//mux.HandleFunc("/wsmock", handlers.WSMockHandler)

	// Health Handler
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Events or UI Handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method, %s", r.Method)
		if r.Method == "POST" {
			t.ServeHTTP(w, r)
			return
		}
		handlers.RootHandler(w, r)
	})

	a := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(a, mux))

}
