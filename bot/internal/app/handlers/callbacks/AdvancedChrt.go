package callbacks

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func AdvancedChrt(
	profile *profile.Service,
	chart *chart.Service,
) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chartActionsData := strings.Split(update.CallbackQuery.Data, ":")
		if len(chartActionsData) > 2 {
			mesId, err := strconv.Atoi(chartActionsData[3])
			if err == nil {
				params := &bot.DeleteMessageParams{
					ChatID: update.CallbackQuery.From.ID,
				}
				for i := update.CallbackQuery.Message.Message.ID - 1; i != mesId-1; i-- {
					params.MessageID = i
					b.DeleteMessage(ctx, params)
				}
			}
		}
		IdUser := update.CallbackQuery.From.ID

		IdChartStr := strings.Trim(chartActionsData[1], " ")
		keyboard := [][]models.InlineKeyboardButton{
			{{Text: "🕯 Интерпритация", CallbackData: fmt.Sprintf("InterpritationChart: %v", IdChartStr)}},
			{{Text: "🔮 Прогноз", CallbackData: fmt.Sprintf("CreateForecast:%v", IdChartStr)}},
			{{Text: "🗑 Удалить карту", CallbackData: fmt.Sprintf("RemoveChart:%v", IdChartStr)}},
			{{Text: "👈 Вернуться к картам", CallbackData: "charts"}},
		}

		params := &bot.EditMessageTextParams{
			ChatID:    IdUser,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text:      "📜 Информация о выбранной натальной карте",
			ParseMode: models.ParseModeHTML,
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		}
		_, err := b.EditMessageText(ctx, params)
		if err != nil && strings.Contains(err.Error(), "Forbidden") {
			profile.DeleteProfile(ctx, IdUser)
		}
	}
}
