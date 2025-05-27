package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
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
	_, err = kadm.NewClient(client).CreateTopic(context.Background(), 1, -1, nil, topic)
	if err != nil {
		//panic(err)
		fmt.Println(err)
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
