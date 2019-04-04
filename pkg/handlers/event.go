package handlers

import (
	"context"
	"log"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/mchmarny/tevents/pkg/twitter"
)

const (
	knownPublisherTokenName = "token"
)

// TwitterReceiver later
type TwitterReceiver struct {
}

// Receive receives
func (t *TwitterReceiver) Receive(ctx context.Context, event cloudevents.Event, _ *cloudevents.EventResponse) error {

	// TODO: Implement in source
	// If token knownPublisherToken
	//// check for presence of publisher token
	// var srcToken string
	// ctx := event.Context.AsV02()
	// if ctx.Extensions != nil {
	// 	if t, ok := ctx.Extensions[knownPublisherTokenName]; ok {
	// 		if srcToken, ok = t.(string); !ok {
	// 			log.Printf("Invalid request (%s missing)", knownPublisherTokenName)
	// 			return
	// 		}
	// 	}
	// }

	// check validity of poster token
	// if srcToken == "" {
	// 	log.Printf("nil token: %s", srcToken)
	// 	return
	// } else if knownPublisherToken != srcToken {
	// 	log.Printf("invalid token: %s", srcToken)
	// 	return
	// }

	log.Printf("Event: %v", event)

	data := twitter.SimpleTweet{}
	if err := event.DataAs(&data); err != nil {
		// the content is not a Twitter message
		log.Printf("Failed to DataAs: %s", err.Error())
		return err
	}

	manager.broadcast <- data

	return nil

}
