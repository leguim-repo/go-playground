package publisher

import (
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker   = "tcp://localhost:1883"
	clientID = "go-mqtt-client"
	topic    = "/test/message"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func generateRandomMessage(messages []string) string {
	return messages[rand.Intn(len(messages))]
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RickAndRoll(randomMessage bool) {
	var message string
	var messageIndex int = 0

	rickAndRollMessages := []string{
		"1 - Never gonna give you up, never gonna let you down",
		"2 - Never gonna run around and desert you",
		"3 - Never gonna make you cry, never gonna say goodbye",
		"4 - Never gonna tell a lie and hurt you",
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		if randomMessage {
			message = generateRandomMessage(rickAndRollMessages)
		} else {
			//fmt.Println("messageIndex:", messageIndex, " message: ", rickAndRollMessages[messageIndex])
			message = rickAndRollMessages[messageIndex]
			messageIndex = messageIndex + 1
			if messageIndex > len(rickAndRollMessages)-1 {
				messageIndex = 0
			}
		}

		token := client.Publish(topic, 0, false, message)
		token.Wait()
		fmt.Printf("Published message: %s\n", message)
		time.Sleep(100 * time.Millisecond)

	}
}
