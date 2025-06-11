package commands

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Izumra/Magistus/bot/internal/services/profile"
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

		photo, err := os.ReadFile("bot/assets/mona_greeting.jpg")
		if err != nil {
			logger.Info("Photo wasn't opened", slog.Any("err", err))
			return
		}

		keyboard := [][]models.InlineKeyboardButton{
			{{Text: "Создать карту", CallbackData: "createchart"}},
			{{Text: "Мои карты", CallbackData: "charts"}},
		}
		params := &bot.SendPhotoParams{
			ChatID: update.Message.Chat.ID,
			Photo: &models.InputFileUpload{
				Filename: "mona_greeting.jpg",
				Data:     bytes.NewReader(photo),
			},
			Caption: strings.ReplaceAll(
				`Привет!
			Меня зовут Мона Мегистус, но друзья в основном обращаются ко мне просто по имени.
			Я астролог из свободолюбивого города Мондштадт.
			
			Судьба каждого человека предопределена с рождения.🔮
			Ты не поверишь, но 🌌 могут многое расскать о тебе и твоём будущем! 
			Я помогу тебе узнать о себе больше 😽`,
				"\t",
				"",
			),
			ParseMode: models.ParseModeHTML,
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		}

		_, err = b.SendPhoto(ctx, params)
		if err != nil {
			logger.Info("Message was't sended", slog.Any("err", err))
		}
	}
}
