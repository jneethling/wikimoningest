package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	var err error
	var (
		kafkaBroker    = flag.String("kafkaBroker", "kafka-server:9092", "Kafka server URL")
		kafkaWikiTopic = flag.String("kafkaWikiTopic", "wikimon", "Outgoing wikimon topic")
	)
	flag.Parse()

	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Could not initialise Zap logger: %v", err)
	}
	defer logger.Sync()

	kafkaController, err := NewKafkaClient(*kafkaBroker)
	if err != nil {
		log.Fatalf("Could not initialize kafka broker: %v", err)
	}

	producer := new(WikimonProducer)
	err = producer.Init(*logger, kafkaController, *kafkaWikiTopic)
	if err != nil {
		log.Fatalf("Could not initialize kafka producer: %v", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		osCall := <-interrupt
		log.Print("System call: " + osCall.String())
		cancel()
	}()

	data := make(chan []byte)
	go producer.ProduceMsg(data)

	err = ConnectWS(ctx, logger, data)
	if err != nil {
		log.Fatal(err)
	}

	close(data)
	log.Print("App stopped by system call")
}

func ConnectWS(ctx context.Context, logger *zap.Logger, data chan<- []byte) error {
	u := url.URL{Scheme: "ws", Host: "wikimon.hatnote.com:9000"}
	logger.Info("Connecting to " + u.String())
	con, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Error("Handshake failed with status " + strconv.Itoa(resp.StatusCode))
		return err
	}
	go func() {
		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				logger.Error("Read: " + err.Error())
				return
			}
			data <- message
		}
	}()
	<-ctx.Done()
	return nil
}

// NewKafkaClient creates a synchronized producer
func NewKafkaClient(broker string) (sarama.SyncProducer, error) {

	brokers := []string{broker}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	return sarama.NewSyncProducer(brokers, config)
}
