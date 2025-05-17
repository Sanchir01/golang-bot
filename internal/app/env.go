package app

import (
	"github.com/go-telegram/bot"
)

type Env struct {
	Reader *KafkaConsumer
	Bot    *bot.Bot
}

func NewENV() (*Env, error) {

	opts := []bot.Option{}

	b, err := bot.New("7829990527:AAHaVWh16TQoNI7AiYPR-VwA-jFc-PxtvwA", opts...)
	if err != nil {
		return nil, err
	}
	kafkaReader, err := NewConsumer("test-topic", "localhost:9092", b)
	if err != nil {
		return nil, err
	}
	return &Env{
		Reader: kafkaReader,
		Bot:    b,
	}, nil
}
