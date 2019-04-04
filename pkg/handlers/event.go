package handlers

import (
	"log"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/mchmarny/tevents/pkg/twitter"
)

const (
	knownPublisherTokenName = "token"
)

// TwitterEventsReceived handles the cloud event post
func TwitterEventsReceived(event cloudevents.Event) {

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
		return
	}

	manager.broadcast <- data

}
