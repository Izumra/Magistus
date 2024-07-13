package callbacks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func UserCharts(
	prof *profile.Service,
	chart *chart.Service,
) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatId := update.CallbackQuery.Message.Message.Chat.ID

		_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    chatId,
			MessageID: update.CallbackQuery.Message.Message.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "Forbidden") {
				prof.DeleteProfile(ctx, chatId)
			}
			return
		}

		errParams := &bot.SendMessageParams{
			ChatID: chatId,
		}

		charts, err := prof.ListCharts(ctx, chatId)
		if err != nil {
			if errors.Is(err, profile.ErrChartsNotFound) {
				errParams.ReplyMarkup = models.InlineKeyboardMarkup{
					InlineKeyboard: [][]models.InlineKeyboardButton{
						{{Text: "🗞 Создать карту", CallbackData: "createchart"}},
					},
				}
			}

			errParams.Text = err.Error()
			_, err := b.SendMessage(ctx, errParams)
			if err != nil && strings.Contains(err.Error(), "Forbidden") {
				prof.DeleteProfile(ctx, chatId)
			}
		}

		var keyboard [][]models.InlineKeyboardButton
		for i := range charts {
			keyboard = append(
				keyboard,
				[]models.InlineKeyboardButton{
					{
						Text:         charts[i].Title,
						CallbackData: fmt.Sprintf("AdvancedChrt:%v", charts[i].Id),
					},
				},
			)
		}

		params := &bot.SendMessageParams{
			ChatID:    chatId,
			Text:      "📜 Список ваших натальных карт получен",
			ParseMode: models.ParseModeHTML,
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		}
		_, err = b.SendMessage(ctx, params)
		if err != nil && strings.Contains(err.Error(), "Forbidden") {
			prof.DeleteProfile(ctx, chatId)
		}
	}
}
