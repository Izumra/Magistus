package callbacks

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SelectChartParam(
	profile *profile.Service,
	chart *chart.Service,
) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		interpritData := strings.Split(update.CallbackQuery.Data, ":")

		chartId, err := strconv.ParseInt(interpritData[1], 10, 0)
		if err != nil {
			log.Println(err)
			return
		}

		paramId, err := strconv.Atoi(interpritData[2])
		if err != nil {
			return
		}

		chart, err := chart.GetChart(ctx, chartId)
		if err != nil {
			return
		}

		exp := regexp.MustCompile(`<h3.*?>(.*?)</h3>`)
		descers := exp.Split(chart.Interpritaion, -1)

		expPTag := regexp.MustCompile(`<p.*?>|</p>`)
		expDivTag := regexp.MustCompile(`<div.*?>|</div>`)
		descers[paramId+1] = expDivTag.ReplaceAllString(descers[paramId+1], "")

		descers[paramId+1] = expPTag.ReplaceAllString(descers[paramId+1], "\n")
		descers[paramId+1] = strings.ReplaceAll(descers[paramId+1], "<p>", "\n")

		messages := []string{}
		if len(descers[paramId+1]) > 3000 {

			words := strings.Split(descers[paramId+1], " ")

			var currentMessage []string
			currentSize := 0

			for _, word := range words {

				if currentSize+len(word) <= 3000 {
					currentMessage = append(currentMessage, word)
					currentSize += len(word)
				} else {
					message := strings.Join(currentMessage, " ")
					messages = append(messages, message)
					currentMessage = []string{word}
					currentSize = len(word)
				}
			}
		}

		if len(messages) == 0 {
			return
		}
		firstMessageParamsViewIterpritation := &bot.EditMessageTextParams{
			ChatID:    update.CallbackQuery.From.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
			ParseMode: "HTML",
			Text:      messages[0],
		}
		mes, err := b.EditMessageText(ctx, firstMessageParamsViewIterpritation)
		if err != nil {
			if strings.Contains(err.Error(), "Forbidden") {
				profile.DeleteProfile(ctx, update.CallbackQuery.From.ID)
			}
		}

		messageParamsViewIterpritation := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: "HTML",
		}
		for i := 1; i < len(messages)-1; i++ {

			messageParamsViewIterpritation.Text = messages[i]
			_, err = b.SendMessage(ctx, messageParamsViewIterpritation)
			if err != nil {
				if strings.Contains(err.Error(), "Forbidden") {
					profile.DeleteProfile(ctx, update.CallbackQuery.From.ID)
				}
				return
			}
		}

		keyboard := make([][]models.InlineKeyboardButton, 1)
		keyboard[0] = []models.InlineKeyboardButton{
			{
				Text:         "👈 Вернуться к карте",
				CallbackData: fmt.Sprintf("AdvancedChrt:%v:deleteTo:%d", chart.Id, mes.ID),
			},
		}
		messageParamsViewIterpritation.Text = messages[len(messages)-1]
		messageParamsViewIterpritation.ReplyMarkup = models.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		}
		_, err = b.SendMessage(ctx, messageParamsViewIterpritation)
		if err != nil {
			if strings.Contains(err.Error(), "Forbidden") {
				profile.DeleteProfile(ctx, update.CallbackQuery.From.ID)
			}
			return
		}

	}
}
