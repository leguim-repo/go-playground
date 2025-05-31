package producer

import (
	"context"
	"flag"
	"fmt"
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

// MiHookCombinado implementa BrokerConnectHook y hooks de producción.
type MiHookCombinado struct{}

// ---- Implementación de kgo.BrokerConnectHook ----

// OnConnect es llamado cuando el cliente se conecta a un broker.
func (h *MiHookCombinado) OnConnect(meta kgo.BrokerMetadata, dialDur time.Duration, conn net.Conn, err error) {
	if err != nil {
		fmt.Printf("HOOK (OnConnect): Error al conectar con broker %s (ID: %d): %v (duración: %s)\n", meta.Host, meta.NodeID, err, dialDur)
		return
	}
	if conn != nil {
		fmt.Printf("HOOK (OnConnect): Conectado a broker %s (ID: %d) en %s. Addr local: %s, Addr remota: %s\n",
			meta.Host, meta.NodeID, dialDur, conn.LocalAddr(), conn.RemoteAddr())
	} else {
		// Esto podría ocurrir si el hook es llamado con un error antes de que la conexión se establezca.
		fmt.Printf("HOOK (OnConnect): Intento de conexión a broker %s (ID: %d) finalizado (duración: %s), pero la conexión es nil (posiblemente debido a un error previo).\n", meta.Host, meta.NodeID, dialDur)
	}
}

// OnDisconnect es llamado cuando el cliente se desconecta de un broker.
func (h *MiHookCombinado) OnDisconnect(meta kgo.BrokerMetadata, conn net.Conn) {
	if conn != nil {
		fmt.Printf("HOOK (OnDisconnect): Desconectado del broker %s (ID: %d). Addr local: %s, Addr remota: %s\n",
			meta.Host, meta.NodeID, conn.LocalAddr(), conn.RemoteAddr())
	} else {
		fmt.Printf("HOOK (OnDisconnect): Desconectado del broker %s (ID: %d) (conn era nil)\n", meta.Host, meta.NodeID)
	}
}

func (h *MiHookCombinado) OnBrokerWrite(meta kgo.BrokerMetadata, key int16, bytesWritten int, writeWait, timeToWrite time.Duration, err error) {
	if err != nil {
		fmt.Printf("HOOK (OnBrokerWrite) fail host: %s err: %s", meta.Host, err)
		return
	}
}
func PlaygroundRedPandaProducerHook() {
	flag.Parse()

	hooks := &MiHookCombinado{}

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
		for range time.Tick(time.Second) {
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
