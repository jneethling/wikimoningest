package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInvalidBroker(t *testing.T) {
	_, err := NewKafkaClient("badBrokerString")
	assert.Error(t, err)
}

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	logger, _ = zap.NewProduction()
	data := make(chan []byte)
	err := ConnectWS(ctx, logger, data)
	assert.NoError(t, err)
}
