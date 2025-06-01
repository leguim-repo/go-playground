package producer

import (
	"context"
	"flag"
	"fmt"
	"go-playground/pkg/thelogger"
	"net"
	_ "net/http/pprof"
	"strconv"
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
	h.logger.Info("HOOK: OnBrokerConnect called")

	if err != nil {
		h.logger.Info("HOOK (OnConnect): Error connect to broker host: " + meta.Host + " (ID: " + strconv.FormatInt(int64(meta.NodeID), 10) + "): error: " + err.Error() + " (initDur: " + initDur.String() + ")")
		return
	}
	if conn != nil {
		h.logger.Info("HOOK (OnConnect): Connected to broker host: " + meta.Host + " (ID: " + strconv.FormatInt(int64(meta.NodeID), 10) + ") initDur: " + initDur.String() + " LocalAddr: " + conn.LocalAddr().String() + " RemoteAddr: " + conn.RemoteAddr().String())
	} else {
		// // This case (err == nil and conn == nil) is unexpected based on typical franz-go behavior, logging as a warning.
		h.logger.Warn("HOOK (OnConnect): Broker connection attempt " + meta.Host + " (ID: " + strconv.FormatInt(int64(meta.NodeID), 10) + ") finished (duration: " + initDur.String() + "), connection is nil).")
	}
}

// OnBrokerDisconnect is called when a connection to a broker is closed.
func (h *BrokerHooks) OnBrokerDisconnect(meta kgo.BrokerMetadata, conn net.Conn) {
	h.logger.Info("HOOK: OnDisconnect called")

	if conn != nil {
		h.logger.Info("HOOK (OnDisconnect): Disconnected from broker " + meta.Host + " (ID: " + strconv.FormatInt(int64(meta.NodeID), 10) + "). LocalAddr: " + conn.LocalAddr().String() + ", RemoteAddr: " + conn.RemoteAddr().String())
	} else {
		h.logger.Warn("HOOK (OnDisconnect): Disconnected from broker " + meta.Host + " (ID: " + strconv.FormatInt(int64(meta.NodeID), 10) + ") (conn = nil)")
	}
}

// OnBrokerWrite is called after write to a broker.
func (h *BrokerHooks) OnBrokerWrite(meta kgo.BrokerMetadata, key int16, bytesWritten int, writeWait, timeToWrite time.Duration, err error) {
	h.logger.Info("HOOK: OnBrokerWrite called")

	if err != nil {
		h.logger.Error("HOOK (OnBrokerWrite): Error when writing in broker " + meta.Host + " (ID: " + strconv.FormatInt(int64(meta.NodeID), 10) + "), key " + strconv.FormatInt(int64(key), 10) + ": err: " + err.Error() + " (writeWait: " + strconv.Itoa(int(writeWait)) + ", timeToWrite: " + strconv.Itoa(int(timeToWrite)) + ")")
	} else {
		h.logger.Info("HOOK (OnBrokerWrite): Successful writing in broker " + meta.Host + " (ID: " + strconv.FormatInt(int64(meta.NodeID), 10) + "), key " + strconv.FormatInt(int64(key), 10) + ". Bytes: " + strconv.Itoa(bytesWritten) + " (writeWait: " + strconv.Itoa(int(writeWait)) + ", timeToWrite: " + strconv.Itoa(int(timeToWrite)) + ")")
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
		kgo.WithHooks(hooks),
		kgo.DefaultProduceTopic(*topic),
		kgo.AllowAutoTopicCreation(),
		// Client logger (kgo.BasicLogger) is commented out; hooks use their own custom logger.
		//kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, func() string {
		//	return time.Now().Format("[2006-01-02 15:04:05.999] ")
		//})),
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
