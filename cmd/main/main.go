package main

import (
	"context"
	"fmt"
	"github.com/Sanchir01/kafka-tg/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	env, err := app.NewENV()
	if err != nil {
		panic(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Signal(syscall.SIGTERM), syscall.SIGINT)
	defer cancel()

	go func() {
		fmt.Println("✅ Kafka consumer started")
		if err := env.Reader.Consume(ctx); err != nil {
			fmt.Printf("Kafka consumer error: %v\n", err)
			cancel()
		}
	}()
	env.Bot.Start(ctx)
}
