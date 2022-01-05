package main

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	wikimon "github.com/jneethling/wikimoncodec/wikimon"
	"go.uber.org/zap"
)

// TransactionProducer holds the kafka settings
type WikimonProducer struct {
	zapLogger       zap.Logger
	wikiCodec       wikimon.AvroCodec
	kafkaController sarama.SyncProducer
	kafkaTopic      string
}

// Init the producer parameters
func (p *WikimonProducer) Init(logger zap.Logger, kafkaController sarama.SyncProducer, kafkaTopic string) error {
	p.zapLogger = logger
	p.kafkaTopic = kafkaTopic
	p.kafkaController = kafkaController

	wikiCodec, err := wikimon.WikiAvroCodec()
	if err != nil {
		p.zapLogger.Error("Error getting AVRO codec")
		return err
	}
	p.wikiCodec = *wikiCodec

	return nil
}

func (p *WikimonProducer) ProduceMsg(data <-chan []byte) {

	for wsmsg := range data {

		var native wikimon.Wikimon
		err := json.Unmarshal(wsmsg, &native)
		if err != nil {
			p.zapLogger.Error("Error converting message to native golang struct: " + err.Error())
			continue
		}

		m, err := json.Marshal(native)
		if err != nil {
			p.zapLogger.Error("Error marshalling message: " + err.Error())
			continue
		}

		var msg = make(map[string]interface{})
		err = json.Unmarshal(m, &msg)
		if err != nil {
			p.zapLogger.Error("Error unmarshalling message to go map string interface: " + err.Error())
			continue
		}

		encodedMsg, err := p.wikiCodec.BinaryFromNative(nil, msg)
		if err != nil {
			p.zapLogger.Error("Error encoding message to AVRO")
			continue
		}
		avroMsg := &sarama.ProducerMessage{Topic: p.kafkaTopic, Value: sarama.StringEncoder(encodedMsg)}
		partition, offset, err := p.kafkaController.SendMessage(avroMsg)
		if err != nil {
			p.zapLogger.Error("Error writing message to kafka topic: " + err.Error())
			continue
		}
		p.zapLogger.Info("Produced message to partition " + fmt.Sprint(partition) + " with offset " + fmt.Sprint(offset))

	}

}
