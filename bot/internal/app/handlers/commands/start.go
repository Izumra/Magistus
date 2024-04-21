package commands

import (
	"context"
	"log/slog"

	"github.com/Izumra/Magistus/bot/internal/services/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(
	logger *slog.Logger,
	profile *profile.Service,
) bot.HandlerFunc {
	op := "command.Start"
	logger = logger.With(slog.String("handler", op))

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		user := update.Message.From

		err := profile.AuthProfile(ctx, user)
		if err != nil {
			return
		}

		keyboard := [][]models.InlineKeyboardButton{
			{{Text: "🗞 Создать карту", CallbackData: "createchart"}},
			{{Text: "📜 Мои карты", CallbackData: "charts"}},
		}
		params := &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "Some greeting text",
			ParseMode: models.ParseModeHTML,
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		}
		m, err := b.SendMessage(ctx, params)
		if err != nil {
			logger.Info("Message was't sended", err, slog.Any("message", m))
		}
	}
}
