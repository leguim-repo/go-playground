package wigeracedata

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go-playground/pkg/datetimeutils"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
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
	EventPid        []int  `json:"eventPid"`
	ClientLocalTime string `json:"clientLocalTime"`
}

var globalMessageCounter int

func incrementGlobalMessageCounter() {
	globalMessageCounter++
}

func getCurrentDirectory() string {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current directory: %s", currentDirectory)
	return currentDirectory
}

func saveData(messageWS string, sessionName string) {

	jsonData := []byte(messageWS)
	incrementGlobalMessageCounter()
	currentTimeUTC := datetimeutils.Now()
	partition, _ := datetimeutils.CreatePartitionStamp(currentTimeUTC)
	formattedTimeStamp, err := datetimeutils.CreateFileTimeStamp(currentTimeUTC)

	dataLakePath := getCurrentDirectory() + "/datalake/" + partition + "/" + sessionName + "/"

	if err := os.MkdirAll(dataLakePath, 0755); err != nil {
		log.Printf("Error creando directorio: %s", err)
		return
	}

	fileName := dataLakePath + sessionName + "_" + strconv.Itoa(globalMessageCounter) + "_" + string(formattedTimeStamp) + ".json"

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing file: %s", err)
		return
	}
	log.Printf("Data saved: %s", fileName)
}

func WigeRaceData() {
	raceDataInitMessage := leaderBoardInitMessage{
		EventId:         the24HNurbugring,
		EventPid:        []int{0, 4},
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
			_, message, err := conn.ReadMessage()
			// Handle received messages and read errors
			if err != nil {
				log.Printf("Error while receiving message: %v", err)
				return
			}
			//log.Printf("Message received (type %d): %s", messageType, string(message))
			saveData(string(message), "24HN_CLASSIC_QUALY_1")
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
