package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/mchmarny/tevents/pkg/handlers"
	"github.com/mchmarny/tevents/pkg/utils"
	cehttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"

	"golang.org/x/net/websocket"
)

func main() {

	ctx := context.Background()
	port := utils.MustGetEnv("PORT", "8080")

	// Configs
	handlers.InitHandlers()

	// Static
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	// UI Handlers
	http.HandleFunc("/", handlers.RootHandler)
	http.Handle("/ws", websocket.Handler(handlers.WSHandler))

	// Health Handler
	http.HandleFunc("/_health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Ingres API Handler
	t, err := cehttp.New(cehttp.WithMethod("POST"), cehttp.WithPath("/v2/twitter"))
	if err != nil {
		log.Fatalf("failed to create cloudevents transport, %s", err.Error())
	}

	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create cloudevents client, %s", err.Error())
	}

	log.Println("Starting twitter receiver...")
	if err := c.StartReceiver(ctx, handlers.TwitterEventsReceived); err != nil {
		log.Fatalf("failed to start twitter events receiver, %s", err.Error())
	}
	log.Printf("Server starting on port %s \n", port)

	// Block until done.
	<-ctx.Done()
}
