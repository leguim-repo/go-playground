package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

func DemoConsumer() {
	seeds := []string{"localhost:19092"}
	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("my-foobar-group"),
		kgo.ConsumeTopics("foobar"),
	)
	if err != nil {
		log.Fatal("Error creating client:", err)
	}
	defer client.Close()

	ctx := context.Background()

	log.Println("Starting message consumption...")

	// Infinite loop
	for {
		fetches := client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("Error consuming messages: %v\n", err)
			}
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			fmt.Printf("Received message from topic %s [%d] offset %d: %s = %s\n",
				record.Topic,
				record.Partition,
				record.Offset,
				string(record.Key),
				string(record.Value),
			)
		}
	}
}
