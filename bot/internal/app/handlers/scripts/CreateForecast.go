package scripts

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Izumra/Magistus/bot/internal/lib/converter"
	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func CreateForecast(
	logger *slog.Logger,
	profile *profile.Service,
	chart *chart.Service,
) bot.HandlerFunc {
	op := "scripts.CreateForecast"
	logger = logger.With(slog.String("handler", op))

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		paramChartID := strings.Split(update.CallbackQuery.Data, ":")[1]
		chartID, err := strconv.ParseInt(paramChartID, 10, 0)
		if err != nil {
			return
		}

		IdUser := update.CallbackQuery.From.ID

		params := &bot.EditMessageTextParams{
			ChatID:    IdUser,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text:      "📝 Отправьте боту координаты места вашего нахождения",
		}
		_, err = b.EditMessageText(ctx, params)
		if err != nil {
			logger.Info("Message was't sended", slog.Any("cause", err))
			return
		}

		//Step 1: getting current user position
		cordsChan := make(chan *models.Location)

		stepHandler1 := b.RegisterHandlerMatchFunc(
			handlerMatchFunc(IdUser),
			func(ctx context.Context, b *bot.Bot, update *models.Update) {
				if update.Message != nil && update.Message.Location != nil {
					cordsChan <- update.Message.Location
					return
				}

				params := &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "✖️ Неверный формат данных локации, попробуйте снова",
				}

				_, err := b.SendMessage(ctx, params)
				if err != nil {
					logger.Info("Message was't sended", slog.Any("cause", err))
				}
				close(cordsChan)
			},
		)

		cords, ok := <-cordsChan
		b.UnregisterHandler(stepHandler1)
		if !ok {
			return
		}
		close(cordsChan)

		mapCords := converter.ConvertCordsToMapCords(cords.Longitude, cords.Latitude)

		chanTimeZone := make(chan string, 1)
		go reqForTheTimeZone(chanTimeZone, cords.Latitude, cords.Longitude)

		timeZone := <-chanTimeZone
		close(chanTimeZone)

		loc, err := time.LoadLocation(timeZone)
		if err != nil {
			return
		}

		timeForecast := time.Now()
		timestamp := time.Date(timeForecast.Year(), timeForecast.Month(), timeForecast.Day(), timeForecast.Hour(), timeForecast.Minute(), timeForecast.Second(), 0, loc)

		forecast, err := chart.CreateForecast(
			ctx,
			chartID,
			timestamp,
			mapCords,
		)
		if err != nil {
			params := &bot.SendMessageParams{
				ChatID: IdUser,
				Text:   "✖️ Бот не смог составить прогноз",
			}
			_, err := b.SendMessage(ctx, params)
			if err != nil && strings.Contains(err.Error(), "Forbidden") {
				profile.DeleteProfile(ctx, IdUser)
			}
		}

		exp := regexp.MustCompile(`<h3(.*?)>(.*?)</h3>`)
		descers := exp.Split(forecast, -1)

		expPTag := regexp.MustCompile(`<p(.*?)>|</p>`)
		expDivTag := regexp.MustCompile(`<div(.*?)>|</div>`)

		descers[1] = expDivTag.ReplaceAllString(descers[1], "")
		descers[1] = expPTag.ReplaceAllString(descers[1], "\n")

		messages := []string{}
		if len(descers[1]) > 3000 {

			words := strings.Split(descers[1], " ")

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

		var sendedMessage *models.Message
		messageParamsViewIterpritation := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: "HTML",
		}
		for i := 0; i < len(messages)-1; i++ {

			messageParamsViewIterpritation.Text = messages[i]
			sendedMessage, err = b.SendMessage(ctx, messageParamsViewIterpritation)
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
				CallbackData: fmt.Sprintf("AdvancedChrt: %v:deleteTo:%d", chartID, sendedMessage.ID),
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
