package main

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	ctx := context.Background()

	seeds := []string{"localhost:19092"}

	client, err := kgo.NewClient(kgo.SeedBrokers(seeds...))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	topic := "foobar"
	// Create a RedPanda topic
	_, err = kadm.NewClient(client).CreateTopic(context.Background(), 1, -1, nil, topic)
	if err != nil {
		//panic(err) // if topic exist trigger panic error
		fmt.Println(err)
	}

	// Getting list of topics. Topics with prefix _ are internal
	currentTopics, err := kadm.NewClient(client).ListTopics(context.Background())
	for _, currentTopic := range currentTopics {
		if !strings.HasPrefix(currentTopic.Topic, "_") {
			fmt.Println("currentTopic: ", currentTopic.Topic)
		}
	}

	count := 5
	wg := sync.WaitGroup{}
	wg.Add(count)

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
