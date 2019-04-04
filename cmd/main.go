package main

import (
	"context"
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

	ctx := context.Background()
	portStr := utils.MustGetEnv("PORT", "8080")

	port, err := strconv.Atoi(portStr)
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

	// UI Handlers
	mux.HandleFunc("/", handlers.RootHandler)
	mux.Handle("/ws", websocket.Handler(handlers.WSHandler))

	// Health Handler
	mux.HandleFunc("/_health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Ingres API Handler
	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithMethod("POST"),
		cloudevents.WithPath("/v2/twitter"),
		cloudevents.WithPort(port),
	)
	if err != nil {
		log.Fatalf("failed to create cloudevents transport, %s", err.Error())
	}
	// Provide extra handlers
	t.Handler = mux

	c, err := cloudevents.NewClient(t, cloudevents.WithUUIDs(), cloudevents.WithTimeNow())
	if err != nil {
		log.Fatalf("failed to create cloudevents client, %s", err.Error())
	}

	log.Println("Starting twitter receiver...")
	log.Printf("Server starting on port %s \n", port)
	if err := c.StartReceiver(ctx, handlers.TwitterEventsReceived); err != nil {
		log.Fatalf("failed to start twitter events receiver, %s", err.Error())
	}

	// Block until done.
	<-ctx.Done()
}
