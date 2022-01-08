package main

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestProducerOK(t *testing.T) {
	logger, _ = zap.NewProduction()
	config := sarama.NewConfig()
	kafkaController := mocks.NewSyncProducer(&testing.T{}, config)
	kafkaController.ExpectSendMessageAndSucceed()

	producer := new(WikimonProducer)
	err := producer.Init(*logger, kafkaController, "testTopic")
	assert.NoError(t, err)

	var wikimon = []byte(`
	{
		"action": "edit",
		"change_size": -12,
		"flags": null,
		"geo_ip": {
		  "city": "Salisbury",
		  "country_name": "United States",
		  "latitude": 38.3761,
		  "longitude": -75.6086,
		  "region_name": "Maryland"
		},
		"hashtags": [],
		"is_anon": true,
		"is_bot": false,
		"is_minor": false,
		"is_new": false,
		"is_unpatrolled": false,
		"mentions": [],
		"ns": "Main",
		"page_title": "Evanescence (Evanescence album)",
		"parent_rev_id": "775894800",
		"rev_id": "774995266",
		"summary": "/* Credits and personnel */ \"Personnel\" is sufficient",
		"url": "https://en.wikipedia.org/w/index.php?diff=775894800&oldid=774995266",
		"user": "71.200.123.192"
	  }
	`)

	data := make(chan []byte)
	defer close(data)
	go producer.ProduceMsg(data)
	data <- wikimon
	time.Sleep(5 * time.Second)
	assert.Equal(t, 0, producer.Failures)
	assert.Equal(t, 1, producer.Successes)
}

func TestProducerFail(t *testing.T) {
	logger, _ = zap.NewProduction()
	config := sarama.NewConfig()
	kafkaController := mocks.NewSyncProducer(&testing.T{}, config)
	kafkaController.ExpectSendMessageAndFail(sarama.ErrUnknown)

	producer := new(WikimonProducer)
	err := producer.Init(*logger, kafkaController, "testTopic")
	assert.NoError(t, err)

	var wikimon = []byte(`{"invalid": "edit"}`)

	data := make(chan []byte)
	defer close(data)
	go producer.ProduceMsg(data)
	data <- wikimon
	time.Sleep(5 * time.Second)
	assert.Equal(t, 1, producer.Failures)
	assert.Equal(t, 0, producer.Successes)
}
