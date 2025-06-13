package ordergrpc

import (
	"bytes"
	"context"
	"fmt"
	orderv1 "github.com/Sanchir01/auth-proto/gen/go/order"
	"github.com/Sanchir01/kafka-tg/internal/domain"
	"github.com/Sanchir01/kafka-tg/pkg/excel"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
)

type ServeApi struct {
	orderv1.UnimplementedOrderServer
	bot *bot.Bot
	log *slog.Logger
}

func RegisterServer(g *grpc.Server, b *bot.Bot) {
	orderv1.RegisterOrderServer(g, ServeApi{bot: b})
}
func (s ServeApi) TelegramOrder(ctx context.Context, request *orderv1.OrderRequest) (*orderv1.OrderResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in TelegramOrder", "error", r, "stack", string(debug.Stack()))
		}
	}()
	if len(request.Orders) == 0 {
		return nil, status.Error(codes.InvalidArgument, "orders is empty")
	}

	slog.Error("TelegramOrder", "request", request.Orders)
	data := MapOrderMessagesToProducts(request)
	if err := excel.CreateFileExcel(data, "candles.xlsx", "Sheet1"); err != nil {
		return nil, fmt.Errorf("failed to create Excel file: %v", err)
	}

	fileData, errReadFile := os.ReadFile("candles.xlsx")
	if errReadFile != nil {
		fmt.Printf("error read file, %v\n", errReadFile)
		return nil, fmt.Errorf("failed to create Excel file: %v", errReadFile)
	}
	_, err := s.bot.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:   1195173283,
		Document: &models.InputFileUpload{Filename: "Заказ", Data: bytes.NewReader(fileData)},
		Caption:  "Заказ",
	})
	if err != nil {
		log.Printf("❗️ Telegram error: %v", err)
		return nil, err
	}
	return &orderv1.OrderResponse{}, nil
}

func MapOrderMessagesToProducts(orders *orderv1.OrderRequest) []domain.ProductWithQuantity {
	products := make([]domain.ProductWithQuantity, 0, len(orders.Orders))

	for _, order := range orders.Orders {

		product := domain.ProductWithQuantity{
			Candles: domain.Candles{
				Title: order.Orders.Title,
				Slug:  order.Orders.Slug,
				Price: order.Orders.Price,
			},
			Quantity: order.Quantity,
		}
		products = append(products, product)
	}

	return products
}
