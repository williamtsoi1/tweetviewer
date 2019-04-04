package handlers

import (
	"log"
	"time"
	"fmt"
	"net/http"

	"github.com/mchmarny/tevents/pkg/twitter"
	"github.com/mchmarny/tevents/pkg/utils"
	"golang.org/x/net/websocket"
)

var (
	manager clientManager
)

type clientManager struct {
    clients    map[*client]bool
    broadcast  chan interface{}
    register   chan *client
    unregister chan *client
}

func (manager *clientManager) start() {
    for {
        select {
        case conn := <-manager.register:
            manager.clients[conn] = true
        case conn := <-manager.unregister:
            if _, ok := manager.clients[conn]; ok {
                close(conn.send)
                delete(manager.clients, conn)
            }
		case message := <-manager.broadcast:
			log.Printf("Broadasting to %d clients: %+v", len(manager.clients), message)
            for conn := range manager.clients {
				log.Printf("Broadasting message to client %s", conn.id)
                select {
                case conn.send <- message:
				default:
                    close(conn.send)
                    delete(manager.clients, conn)
                }
			}
		}
    }
}

func (manager *clientManager) send(message interface{}) {
    for conn := range manager.clients {
		conn.send <- message
    }
}

type client struct {
    id     string
    socket *websocket.Conn
    send   chan interface{}
}

func (c *client) write() {

    defer func() {
        c.socket.Close()
    }()

    for {
        select {
		case message, ok := <-c.send:
			log.Printf("On send %v - %v", message, ok)
            if !ok {
				log.Println("Unable to sent")
                return
            }
			websocket.JSON.Send(c.socket, message)
        }
    }
}


func initWS(){
	manager = clientManager{
		broadcast:  make(chan interface{}, 100),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
	go manager.start()
}

// WSHandler provides backing service for the UI
func WSHandler(ws *websocket.Conn) {
	log.Println("WS connection...")

    client := &client{
		id: utils.MakeUUID(),
		socket: ws,
		send: make(chan interface{}),
	}

    manager.register <- client

    client.write()
}

// WSMockHandler gens mocked tweets
func WSMockHandler(w http.ResponseWriter, r *http.Request) {
	go mockTweets()
	fmt.Fprint(w, "ok")
}


func mockTweets() {
	log.Println("Mocking tweets...")
	for i := 0; i < 100; i++ {
		manager.broadcast <- makeMokeTweet(i)
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
