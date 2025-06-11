package scripts

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Izumra/Magistus/bot/internal/lib/converter"
	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
)

func CreateChart(
	logger *slog.Logger,
	profile *profile.Service,
	chart *chart.Service,
) bot.HandlerFunc {
	op := "scripts.CreateChart"
	logger = logger.With(slog.String("handler", op))

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatId := update.CallbackQuery.Message.Message.Chat.ID

		_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    chatId,
			MessageID: update.CallbackQuery.Message.Message.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "Forbidden") {
				profile.DeleteProfile(ctx, chatId)
			}
			return
		}

		params := &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "📝 Отправьте боту название натальной карты",
		}
		_, err = b.SendMessage(ctx, params)
		if err != nil {
			logger.Info("Message was't sended", slog.Any("cause", err))
			return
		}

		// Step 1: request for the name of the chart
		nameChan := make(chan string)

		stepHandler1 := b.RegisterHandlerMatchFunc(
			handlerMatchFunc(chatId),
			func(ctx context.Context, b *bot.Bot, update *models.Update) {
				if update.Message != nil && update.Message.Text != "" {
					nameChan <- update.Message.Text
					return
				}

				params := &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "✖️ Бот не смог определить имя натальной карты, попробуйте снова",
				}
				_, err := b.SendMessage(ctx, params)
				if err != nil {
					logger.Info("Message was't sended", slog.Any("cause", err))
				}
				close(nameChan)
			},
		)

		nameChart, ok := <-nameChan
		b.UnregisterHandler(stepHandler1)
		if !ok {
			return
		}
		close(nameChan)

		mesBeforeStep2 := &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "📝 Отправьте боту координаты места вашего рождения для построения натальной карты",
		}
		_, err = b.SendMessage(ctx, mesBeforeStep2)
		if err != nil {
			logger.Info("Message was't sended", slog.Any("cause", err))
			return
		}

		// Step 2: request for the coords of the born place
		cordsChan := make(chan *models.Location)

		stepHandler2 := b.RegisterHandlerMatchFunc(
			handlerMatchFunc(chatId),
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
		b.UnregisterHandler(stepHandler2)
		if !ok {
			return
		}
		close(cordsChan)

		mapCords := converter.ConvertCordsToMapCords(cords.Longitude, cords.Latitude)

		mesBeforeStep3 := &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "📝 Отправьте боту дату своего рождения в формате YYYY:MM:DD HH:MM:SS",
		}
		_, err = b.SendMessage(ctx, mesBeforeStep3)
		if err != nil {
			logger.Info("Message was't sended", slog.Any("cause", err))
			return
		}

		chanTimeZone := make(chan string, 1)
		go reqForTheTimeZone(chanTimeZone, cords.Latitude, cords.Longitude)

		// Step 3: request for the date of the born
		bornDateChan := make(chan time.Time)

		stepHandler3 := b.RegisterHandlerMatchFunc(
			handlerMatchFunc(chatId),
			func(ctx context.Context, b *bot.Bot, update *models.Update) {
				paramsBadDataResp := &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "✖️ Неверный формат данных даты рождения, попробуйте снова",
				}

				if update.Message != nil && update.Message.Text != "" {
					userBornDate := update.Message.Text

					timeZone, ok := <-chanTimeZone
					if !ok {
						_, err := b.SendMessage(ctx, paramsBadDataResp)
						if err != nil {
							logger.Info("Message was't sended", slog.Any("cause", err))
						}
						close(bornDateChan)
						return
					}

					loc, err := time.LoadLocation(timeZone)
					if err != nil {
						_, err := b.SendMessage(ctx, paramsBadDataResp)
						if err != nil {
							logger.Info("Message was't sended", slog.Any("cause", err))
						}
						close(bornDateChan)
						return
					}

					dateParts := strings.Fields(userBornDate)
					if len(dateParts) != 2 {
						_, err := b.SendMessage(ctx, paramsBadDataResp)
						if err != nil {
							logger.Info("Message was't sended", slog.Any("cause", err))
						}
						close(bornDateChan)
						return
					}

					dateElems := strings.Split(dateParts[0]+":"+dateParts[1], ":")
					if len(dateElems) != 6 {
						_, err := b.SendMessage(ctx, paramsBadDataResp)
						if err != nil {
							logger.Info("Message was't sended", slog.Any("cause", err))
						}
						close(bornDateChan)
						return
					}

					year, errYconv := strconv.Atoi(dateElems[0])
					month, errMonconv := strconv.Atoi(dateElems[1])
					day, errDconv := strconv.Atoi(dateElems[2])
					hour, errHconv := strconv.Atoi(dateElems[3])
					min, errMinconv := strconv.Atoi(dateElems[4])
					sec, errSecconv := strconv.Atoi(dateElems[5])
					if errYconv != nil || errMonconv != nil || errDconv != nil || errHconv != nil ||
						errMinconv != nil ||
						errSecconv != nil {
						_, err := b.SendMessage(ctx, paramsBadDataResp)
						if err != nil {
							logger.Info("Message was't sended", slog.Any("cause", err))
						}
						close(bornDateChan)
						return
					}

					bornDate := time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
					bornDateChan <- bornDate
					return
				}

				_, err := b.SendMessage(ctx, paramsBadDataResp)
				if err != nil {
					logger.Info("Message was't sended", slog.Any("cause", err))
				}
				close(bornDateChan)
			},
		)

		bornDate, ok := <-bornDateChan
		b.UnregisterHandler(stepHandler3)
		if !ok {
			return
		}

		close(bornDateChan)

		_, err = chart.CreateChart(
			ctx,
			chatId,
			nameChart,
			bornDate,
			mapCords,
		)
		if err != nil {
			params := &bot.SendMessageParams{
				ChatID: chatId,
				Text:   "✖️ Бот не смог построить натальную карту",
			}
			_, err := b.SendMessage(ctx, params)
			if err != nil && strings.Contains(err.Error(), "Forbidden") {
				profile.DeleteProfile(ctx, chatId)
			}
		}

		keyboard := [][]models.InlineKeyboardButton{
			{{Text: "Мои карты", CallbackData: "charts"}},
		}
		respParams := &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "Натальная карта успешно построена",
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		}
		_, err = b.SendMessage(ctx, respParams)
		if err != nil && strings.Contains(err.Error(), "Forbidden") {
			profile.DeleteProfile(ctx, chatId)
		}
	}
}
