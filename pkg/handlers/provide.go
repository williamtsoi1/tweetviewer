package handlers

import (
	"log"
	"time"
	"fmt"

	"github.com/mchmarny/tevents/pkg/twitter"
	"github.com/mchmarny/tevents/pkg/utils"
	"golang.org/x/net/websocket"
)

var (
	eventChannel = make(chan interface{}, 100)
	connections []*websocket.Conn
)

// WSHandler provides backing service for the UI
func WSHandler(ws *websocket.Conn) {
	log.Println("WS connection...")

	connections = append(connections, ws)

	mock := utils.MustGetEnv("MOCK_TWEETS", "no")
	if mock == "yes" {
		go mockTweets()
	}

	for {
		select {
		case m := <-eventChannel:
			for _, w := range connections {
				if err := websocket.JSON.Send(w, m); err != nil {
					log.Printf("Error on write message: %v", err)
				}
			}
		}
	}
}


func mockTweets() {
	for i := 0; i < 100; i++ {
		data := makeMokeTweet(i)
		eventChannel <- data
		time.Sleep(1 * time.Second)
	}
}


func makeMokeTweet(i int) *twitter.SimpleTweet {

	m := "Lorem ipsum dolor sit amet, porttitor turpis mollis, integer ipsum mattis scelerisque aliquam. In volutpat per."

	data := &twitter.SimpleTweet{
		CreatedAt: time.Now().String(),
		IDStr:     fmt.Sprintf("id-%d", i),
		Text: fmt.Sprintf("%s %d", m, i),
		User:      &twitter.SimpleTwitterUser{
			ProfileImageURL: "https://pbs.twimg.com/profile_images/1044758200048898048/dgHKQIOQ_400x400.jpg",
			ScreenName: fmt.Sprintf("@username-%d", i),
		},
	}

	return data

}
