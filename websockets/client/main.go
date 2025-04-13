package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

// Define our message object
type Message struct {
	Type       string        `json:"type"`
	ProductIds []string      `json:"product_ids"`
	Channels   []interface{} `json:"channels"`
}

func main() {
	// Connect to the WebSocket server
	u := url.URL{Scheme: "wss", Host: "ws-feed.exchange.coinbase.com"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Define the ticker channel
	ticker := struct {
		Name       string   `json:"name"`
		ProductIds []string `json:"product_ids"`
	}{
		Name:       "ticker",
		ProductIds: []string{"ETH-BTC", "ETH-USD"},
	}

	// Send initial message
	m := Message{
		Type:       "subscribe",
		ProductIds: []string{"ETH-USD", "ETH-EUR"},
		Channels:   []interface{}{"level2", "heartbeat", ticker},
	}

	msg, _ := json.Marshal(m)

	if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println("write:", err)
		return
	}

	// Start listening for incoming messages
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		log.Printf("recv: %s", message)
	}
}
