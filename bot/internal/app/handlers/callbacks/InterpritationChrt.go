package callbacks

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	ErrUnexpected = "✖️ Неизвестная ошибка, попробуйте еще раз"
)

func InterpritaionChrt(
	profile *profile.Service,
	chart *chart.Service,
) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		interpritData := strings.Split(update.CallbackQuery.Data, ":")
		interpritData[1] = strings.Trim(interpritData[1], " ")

		errResponse := &bot.EditMessageTextParams{
			ChatID:    update.CallbackQuery.From.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
		}

		IdChart, err := strconv.ParseInt(interpritData[1], 10, 0)
		if err != nil {
			errResponse.Text = ErrUnexpected
			_, err = b.EditMessageText(ctx, errResponse)
			if err != nil {
				if strings.Contains(err.Error(), "Forbidden") {
					profile.DeleteProfile(ctx, update.CallbackQuery.From.ID)
				}
			}

			return
		}

		chart, err := chart.GetChart(ctx, IdChart)
		if err != nil {
			errResponse.Text = err.Error()
			_, err = b.EditMessageText(ctx, errResponse)
			if err != nil {
				if strings.Contains(err.Error(), "Forbidden") {
					profile.DeleteProfile(ctx, update.CallbackQuery.From.ID)
				}
			}

			return
		}

		exp := regexp.MustCompile(`<h3 .*?>(.*?)</h3>`)

		titlesSubmatches := exp.FindAllStringSubmatch(chart.Interpritaion, -1)

		keyboard := make([][]models.InlineKeyboardButton, len(titlesSubmatches))
		for i := range titlesSubmatches {
			keyboard[i] = []models.InlineKeyboardButton{
				{
					Text:         titlesSubmatches[i][1],
					CallbackData: fmt.Sprintf("SelectChartParam:%d:%d", IdChart, i),
				},
			}
		}

		paramsViewIterpritation := &bot.EditMessageTextParams{
			ChatID:    update.CallbackQuery.From.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text:      "Выберите параметр интерпритации",
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		}
		_, err = b.EditMessageText(ctx, paramsViewIterpritation)
		if err != nil {
			if strings.Contains(err.Error(), "Forbidden") {
				profile.DeleteProfile(ctx, update.CallbackQuery.From.ID)
			}
		}

	}
}
