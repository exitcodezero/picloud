package socket

import (
	"github.com/exitcodezero/picloud/hub"
	"github.com/exitcodezero/picloud/message"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func writeSocket(socket *websocket.Conn, c *hub.Connection) {
	for {
		m := <-c.Out
		socket.WriteJSON(&m)
	}
}

// Handler handles websocket connections at /connect
func Handler(w http.ResponseWriter, r *http.Request) {

	// Get the "clientName" of the connection from a query parameter
	queryParams := r.URL.Query()
	if len(queryParams["clientName"]) < 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Upgrade the request
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	// Create a Connection instance
	c := hub.NewConnection(socket.RemoteAddr().String(), queryParams["clientName"][0])
	hub.Manager.RegisterConnection(&c)
	defer hub.Manager.Cleanup(&c)

	// Start pushing outbound messages from a goroutine
	go writeSocket(socket, &c)

	// Handle inbound messages
	for {
		m := message.SocketMessage{}
		m.CreatedAt = time.Now().UTC()

		err = socket.ReadJSON(&m)
		if err != nil {
			break
		}

		switch m.Action {
		case "publish":
			hub.Manager.Publish(m)
		case "subscribe":
			hub.Manager.Subscribe(m.Event, &c)
		case "unsubscribe":
			hub.Manager.Unsubscribe(m.Event, &c)
		case "unsubscribe:all":
			hub.Manager.UnsubscribeAll(&c)
		}
	}
}
