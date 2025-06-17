package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type ConfigRedPanda struct {
	SeedBrokers   []string
	Topic         string
	ConsumerGroup string
}

func NewConfigRedPandas() *ConfigRedPanda {
	return &ConfigRedPanda{
		SeedBrokers:   []string{"localhost:19092"},
		Topic:         "foobar",
		ConsumerGroup: "my-foobar-group",
	}
}
func DemoConsumer() {
	configRedPanda := NewConfigRedPandas()

	client, err := kgo.NewClient(
		kgo.SeedBrokers(configRedPanda.SeedBrokers...),
		kgo.ConsumerGroup(configRedPanda.ConsumerGroup),
		kgo.ConsumeTopics(configRedPanda.Topic),
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
