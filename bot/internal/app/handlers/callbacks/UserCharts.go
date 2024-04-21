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
		IdUser := update.CallbackQuery.From.ID

		errParams := &bot.EditMessageTextParams{
			ChatID:    IdUser,
			MessageID: update.CallbackQuery.Message.Message.ID,
		}

		charts, err := prof.ListCharts(ctx, IdUser)
		if err != nil {
			if errors.Is(err, profile.ErrChartsNotFound) {
				errParams.ReplyMarkup = models.InlineKeyboardMarkup{
					InlineKeyboard: [][]models.InlineKeyboardButton{
						{{Text: "🗞 Создать карту", CallbackData: "createchart"}},
					},
				}
			}

			errParams.Text = err.Error()
			_, err := b.EditMessageText(ctx, errParams)
			if err != nil && strings.Contains(err.Error(), "Forbidden") {
				prof.DeleteProfile(ctx, IdUser)
			}
		}

		var keyboard [][]models.InlineKeyboardButton
		for i := range charts {
			keyboard = append(
				keyboard,
				[]models.InlineKeyboardButton{
					{
						Text:         charts[i].Title,
						CallbackData: fmt.Sprintf("AdvancedChrt: %v:deleteTo:0", charts[i].Id),
					},
				},
			)
		}

		params := &bot.EditMessageTextParams{
			ChatID:    IdUser,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text:      "📜 Список ваших натальных карт получен",
			ParseMode: models.ParseModeHTML,
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		}
		_, err = b.EditMessageText(ctx, params)
		if err != nil && strings.Contains(err.Error(), "Forbidden") {
			prof.DeleteProfile(ctx, IdUser)
		}
	}
}
