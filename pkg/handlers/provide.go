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
)


type hub struct {
	clients          map[string]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan interface{}
}

// run the hub
func (h *hub) run() {
	for {
		select {
		case conn := <-h.addClientChan:
			h.addClient(conn)
		case conn := <-h.removeClientChan:
			h.removeClient(conn)
		case m := <-h.broadcastChan:
			h.broadcastMessage(m)
		}
	}
}

// removeClient removes a conn from the pool
func (h *hub) removeClient(conn *websocket.Conn) {
	c := conn.LocalAddr().String()
	log.Printf("Removing client %s: ", c)
	delete(h.clients, c)
}

// addClient adds a conn to the pool
func (h *hub) addClient(conn *websocket.Conn) {
	c := conn.RemoteAddr().String()
	log.Printf("Adding client %s: ", c)
	h.clients[c] = conn
}

// broadcastMessage sends a message to all client conns in the pool
func (h *hub) broadcastMessage(m interface{}) {
	for _, conn := range h.clients {
		err := websocket.JSON.Send(conn, m)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
			h.removeClientChan <- conn
			return
		}
	}
}

// WSHandler provides backing service for the UI
func WSHandler(ws *websocket.Conn) {
	log.Println("WS connection...")

	h := &hub{
		clients:          make(map[string]*websocket.Conn),
		addClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		broadcastChan:    make(chan interface{}),
	}

	go h.run()

	h.addClientChan <- ws

	// Mock - Remove
	mock := utils.MustGetEnv("MOCK_TWEETS", "no")
	log.Printf("Mocking: %s", mock)
	if mock == "yes" {
		go mockTweets()
	}

	for {
		select {
		case m := <-eventChannel:
			h.broadcastChan <- m
		}
	}

}




func mockTweets() {
	log.Println("Mocking tweets...")
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
