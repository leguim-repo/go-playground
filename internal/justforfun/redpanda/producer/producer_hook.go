package producer

import (
	"context"
	"flag"
	"fmt"
	"go-playground/pkg/thelogger"
	"net"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	seedBrokers = flag.String("brokers", "localhost:19092", "comma delimited list of seed brokers")
	topic       = flag.String("topic", "foobar", "topic to consume for metric incrementing")
	produce     = flag.Bool("produce", true, "if true, rather than consume, produce to the topic once per second (value \"foobar\")")
)

// BrokerHooks implements BrokerConnectHook and other hooks of franz-go module
type BrokerHooks struct {
	logger *thelogger.TheLogger
}

// Verify in compilation time if my hooks are implemented correctly
// https://pkg.go.dev/github.com/twmb/franz-go/pkg/kgo@v1.19.4#HookBrokerConnect
var _ kgo.HookBrokerConnect = (*BrokerHooks)(nil)
var _ kgo.HookBrokerDisconnect = (*BrokerHooks)(nil)
var _ kgo.HookBrokerWrite = (*BrokerHooks)(nil)

// OnBrokerConnect is called after a connection to a broker is opened.
func (h *BrokerHooks) OnBrokerConnect(meta kgo.BrokerMetadata, initDur time.Duration, conn net.Conn, err error) {
	if err != nil {
		fmt.Printf("HOOK (OnConnect): Error connect to broker host: %s (ID: %d): error: %v (initDur: %s)\n", meta.Host, meta.NodeID, err, initDur)
		return
	}
	if conn != nil {
		fmt.Printf("HOOK (OnConnect): Connected to broker host: %s (ID: %d) initDur: %s LocalAddr: %s RemoteAddr: %s\n",
			meta.Host, meta.NodeID, initDur, conn.LocalAddr(), conn.RemoteAddr())
	} else {
		// This could happen if the hook is called with an error before the connection is established.
		fmt.Printf("HOOK (OnConnect): Broker connection attempt %s (ID: %d) finished (duration: %s), connection is nil).\n", meta.Host, meta.NodeID, initDur)
	}
}

// OnBrokerDisconnect is called when a connection to a broker is closed.
func (h *BrokerHooks) OnBrokerDisconnect(meta kgo.BrokerMetadata, conn net.Conn) {
	h.logger.Info("HOOK: OnDisconnect called")
	if conn != nil {
		fmt.Printf("HOOK (OnDisconnect): Disconnected from broker %s (ID: %d). LocalAddr: %s, RemoteAddr: %s\n",
			meta.Host, meta.NodeID, conn.LocalAddr(), conn.RemoteAddr())
	} else {
		fmt.Printf("HOOK (OnDisconnect): Disconnected from broker %s (ID: %d) (conn = nil)\n", meta.Host, meta.NodeID)
	}
}

// OnBrokerWrite is called after write to a broker.
func (h *BrokerHooks) OnBrokerWrite(meta kgo.BrokerMetadata, key int16, bytesWritten int, writeWait, timeToWrite time.Duration, err error) {
	if err != nil {
		fmt.Printf("HOOK (OnBrokerWrite): Error when writing in broker %s (ID: %d), key %d: err: %v (writeWait: %s, timeToWrite: %s)\n", meta.Host, meta.NodeID, key, err, writeWait, timeToWrite)
	} else {
		fmt.Printf("HOOK (OnBrokerWrite): Successful writing in broker %s (ID: %d), key %d. Bytes: %d (writeWait: %s, timeToWrite: %s)\n", meta.Host, meta.NodeID, key, bytesWritten, writeWait, timeToWrite)
	}
}

func PlaygroundRedPandaProducerHook() {
	flag.Parse()
	logger := thelogger.NewTheLogger()

	hooks := &BrokerHooks{
		logger: logger,
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(strings.Split(*seedBrokers, ",")...),
		kgo.WithHooks(hooks), // Aquí se registran los hooks
		kgo.DefaultProduceTopic(*topic),
		kgo.AllowAutoTopicCreation(),
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, func() string {
			return time.Now().Format("[2006-01-02 15:04:05.999] ")
		})),
	}
	if !*produce {
		opts = append(opts, kgo.ConsumeTopics(*topic))
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		panic(fmt.Sprintf("unable to create client: %v", err))
	}
	defer cl.Close()

	if *produce {
		for range time.Tick(5 * time.Second) {
			if err := cl.ProduceSync(context.Background(), kgo.StringRecord("foobar")).FirstErr(); err != nil {
				panic(fmt.Sprintf("unable to produce: %v", err))
			}
		}
	} else {
		for {
			cl.PollFetches(context.Background()) // busy work...
		}
	}
}
