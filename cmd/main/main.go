package main

import (
	"context"
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
		env.GRPCSrv.MustStart()
	}()
	env.Bot.Start(ctx)
}
