package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

//NB run docker-compose up to pass this test
func TestInitProducer(t *testing.T) {
	logger, _ = zap.NewProduction()
	kafkaController, err := NewKafkaClient("localhost:9092")
	assert.NoError(t, err)

	producer := new(WikimonProducer)
	err = producer.Init(*logger, kafkaController, "testTopic")
	assert.NoError(t, err)
}
