package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Sanchir01/kafka-tg/internal/domain"
	"github.com/Sanchir01/kafka-tg/pkg/excel"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	bot    *bot.Bot
}

func NewConsumer(topic, broker, groupId string, b *bot.Bot) (*KafkaConsumer, error) {
	cfg := kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupId,
	}
	reader := kafka.NewReader(cfg)

	return &KafkaConsumer{
		reader: reader,
		bot:    b,
	}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context) error {
	defer kc.reader.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Kafka consumer shutdown")
			return nil
		default:
			msg, err := kc.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				log.Printf("❌ fetch error: %v", err)
				continue
			}

			var candles []domain.ProductWithQuantity
			if err := json.Unmarshal(msg.Value, &candles); err != nil {
				log.Printf("⚠️ Failed to unmarshal JSON: %v", err)
				continue
			}
			if err := excel.CreateFileExcel(candles, "candles.xlsx", "Sheet1"); err != nil {
				return fmt.Errorf("failed to create Excel file: %v", err)
			}

			fileData, errReadFile := os.ReadFile("candles.xlsx")
			if errReadFile != nil {
				fmt.Printf("error read file, %v\n", errReadFile)
				return nil
			}
			_, err = kc.bot.SendDocument(ctx, &bot.SendDocumentParams{
				ChatID:   1195173283,
				Document: &models.InputFileUpload{Filename: "Заказ", Data: bytes.NewReader(fileData)},
				Caption:  "Заказ",
			})
			if err != nil {
				log.Printf("❗️ Telegram error: %v", err)
				return nil
			}
			if err := kc.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("⚠️ commit error: %v", err)
			}
		}
	}
}
