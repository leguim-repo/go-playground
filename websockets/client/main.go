package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

// Define our message object
/*
Subscription Message
{
  "type": "subscribe",
  "product_ids": ["ETH-USD", "ETH-EUR"],
  "channels": [
    "heartbeat",
    {
      "name": "ticker",
      "product_ids": ["ETH-BTC", "ETH-USD"]
    }
  ]
}

*/

type subscribeMessage struct {
	Type       string        `json:"type"`
	ProductIds []string      `json:"product_ids"`
	Channels   []interface{} `json:"channels"`
}

func main() {
	// Connect to the WebSocket server
	u := url.URL{Scheme: "wss", Host: "ws-feed.exchange.coinbase.com"}

	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)

	// Define the ticker channel
	tickerSubscribe := struct {
		Name       string   `json:"name"`
		ProductIds []string `json:"product_ids"`
	}{
		Name:       "ticker",
		ProductIds: []string{"ETH-BTC", "ETH-USD"},
	}

	// Send initial message
	m := subscribeMessage{
		Type:       "subscribe",
		ProductIds: []string{"ETH-USD", "ETH-EUR"},
		Channels:   []interface{}{"heartbeat", tickerSubscribe},
	}

	msg, _ := json.Marshal(m)

	if err := connection.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println("write:", err)
		return
	}

	// Start listening for incoming messages
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		log.Printf("recv: %s", message)
	}
}
