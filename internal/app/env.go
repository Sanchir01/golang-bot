package app

import (
	"fmt"
	"log/slog"
	"os"

	cfg "github.com/Sanchir01/kafka-tg/internal/config"
	"github.com/go-telegram/bot"
)

type Env struct {
	Bot     *bot.Bot
	Log     *slog.Logger
	GRPCSrv *App
}

func NewENV() (*Env, error) {
	cfg := cfg.InitConfig()
	fmt.Println("cfg", cfg)
	lg := SetupLogger(cfg.Env)

	opts := []bot.Option{}
	b, err := bot.New(os.Getenv("TOKEN"), opts...)
	if err != nil {
		return nil, err
	}
	gRPCServer := New(lg, ":"+cfg.GRPC.Port, b)
	return &Env{
		GRPCSrv: gRPCServer,
		Log:     lg,
		Bot:     b,
	}, nil
}
