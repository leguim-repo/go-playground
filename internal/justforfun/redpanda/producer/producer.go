package producer

import (
	"context"
	"errors"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"go-playground/pkg/thelogger"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func GetTopics(client *kgo.Client) ([]string, error) {
	var topicsFound []string
	currentTopics, err := kadm.NewClient(client).ListTopics(context.Background())
	if err != nil {
		return topicsFound, errors.New(err.Error())
	}

	for _, currentTopic := range currentTopics {
		if !strings.HasPrefix(currentTopic.Topic, "_") {
			fmt.Println("currentTopic: ", currentTopic.Topic)
			topicsFound = append(topicsFound, currentTopic.Topic)
		}
	}
	return topicsFound, nil
}

func CreateTopic(client *kgo.Client, topic string) error {
	_, err := kadm.NewClient(client).CreateTopic(context.Background(), 1, -1, nil, topic)
	if err != nil {
		return err
	}
	return nil
}

func PlaygroundRedPandaProducer() {
	logger := thelogger.NewTheLogger()
	ctx := context.Background()

	seeds := []string{"localhost:19092"}

	client, err := kgo.NewClient(kgo.SeedBrokers(seeds...))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	topics := []string{"foobar", "foobar2", "foobar3"}
	for _, topic := range topics {
		// Create a RedPanda topic
		logger.Info("Attempt create topic: " + topic)
		err = CreateTopic(client, topic)
		if err != nil {
			logger.Warn(err.Error())
		}
	}
	// Getting list of topics. Topics with prefix _ are internal
	topicsFound, err := GetTopics(client)
	fmt.Println("List of topics found:", topicsFound)
	for _, topic := range topicsFound {
		if !strings.HasPrefix(topic, "_") {
			fmt.Println("topic: ", topic)
		}
	}

	count := 5
	wg := sync.WaitGroup{}
	wg.Add(count)
	topic := "foobar"
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancel()
	for i := range count {
		theMessage := `{"test":"foo","count":` + strconv.Itoa(i) + `,"datetime":"` + time.DateTime + `"}`
		r := &kgo.Record{
			Key:   []byte(strconv.Itoa(i)),
			Topic: topic,
			//Timestamp: time.UnixMilli(int64(i)),
			Value: []byte(theMessage),
		}

		client.Produce(ctx, r, func(_ *kgo.Record, err error) {
			if err != nil {
				panic(err)
			}
			wg.Done()
		})

	}

	wg.Wait()
}
