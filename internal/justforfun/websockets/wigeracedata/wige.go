package wigeracedata

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go-playground/pkg/datetimeutils"
	"log"
	"net/url"
	"os"
	"os/signal"
)

const (
	the24HNurbugring                   = "50"
	fIAEuropeanTruckRacingChampionship = "16"
	nurburgring                        = "21"
	curbStoneTrackEvents               = "111"
	relnoldusLangstreckenCup           = "14"
	vln                                = "20"
	porscheSportsCupDeutschland        = "15"
	nes                                = "112"
	hockenheimRing                     = "41"
	all4Track                          = "113"
	tcrEasternEurope                   = "17"
)

type leaderBoardInitMessage struct {
	EventId         string `json:"eventId"`
	EventPid        string `json:"eventPid"`
	ClientLocalTime string `json:"clientLocalTime"`
}

func WigeRaceData() {
	raceDataInitMessage := leaderBoardInitMessage{
		EventId:         vln,
		EventPid:        "[0, 4]",
		ClientLocalTime: datetimeutils.GetUnixTimestampWithMilliseconds(),
	}

	log.Printf("Wige Race Data WS Client here")

	// Channel for handle (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: "livetiming.azurewebsites.net", Path: "/"}
	log.Printf("Connecting to %s", u.String())

	// Establish connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	// Ensure close connection when finish
	defer conn.Close()

	log.Println("Connected to Wige webSocket server")

	// Send init message
	initMessage, err := json.Marshal(raceDataInitMessage)
	if err != nil {
		panic(err)
	}
	log.Printf("Send init message to %s", string(initMessage))
	err = conn.WriteMessage(websocket.TextMessage, initMessage)
	if err != nil {
		log.Printf("Error while send init message: %v", err)
	} else {
		log.Printf("Init message send ok: %s", string(initMessage))
	}

	// Wait and process messages from server in a goroutine
	done := make(chan struct{})
	// This goroutine listen in loop
	go func() {
		defer close(done)
		// Loop for read messages
		for {
			messageType, message, err := conn.ReadMessage()
			// Handle received messages and read errors
			if err != nil {
				log.Printf("Error while receiving message: %v", err)
				return
			}
			log.Printf("Message received (type %d): %s", messageType, string(message))
		}
	}()

	log.Println("Waiting messages from server... (Press Ctrl+C for exit)")

	// Main loop to keep the connection alive and handle interruptions
	for {
		select {
		case <-done:
			// ...
			log.Println("Closing connection")
			return
		case <-interrupt:
			// ...
			log.Println("Interrupt signal received, exiting...")
			return
		}
	}

}
