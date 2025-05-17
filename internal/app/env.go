package app

import (
	"fmt"
	"os"

	cfg "github.com/Sanchir01/kafka-tg/internal/config"
	"github.com/go-telegram/bot"
)

type Env struct {
	Reader *KafkaConsumer
	Bot    *bot.Bot
}

func NewENV() (*Env, error) {
	con := cfg.InitConfig()

	fmt.Println("cfg", con)
	opts := []bot.Option{}
	b, err := bot.New(os.Getenv("TOKEN"), opts...)
	if err != nil {
		return nil, err
	}
	kafkaReader, err := NewConsumer(con.KafkaConsumer.Consumer.Topic, con.KafkaConsumer.Consumer.Broker[0], con.KafkaConsumer.Consumer.GroupID, b)
	if err != nil {
		return nil, err
	}
	return &Env{
		Reader: kafkaReader,
		Bot:    b,
	}, nil
}
