package app

import (
	"context"
	"log/slog"
	"regexp"
	"syscall"

	"github.com/Izumra/Magistus/bot/internal/app/handlers/callbacks"
	"github.com/Izumra/Magistus/bot/internal/app/handlers/commands"
	"github.com/Izumra/Magistus/bot/internal/app/handlers/scripts"
	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Services struct {
	Chart   *chart.Service
	Profile *profile.Service
}

type Bot struct {
	Instance *bot.Bot
	logger   *slog.Logger
	services *Services
}

func New(
	logger *slog.Logger,
	token string,
	services *Services,
) *Bot {
	bot, err := bot.New(
		token,
		bot.WithSkipGetMe(),
	)
	if err != nil {
		panic(err)
	}

	return &Bot{
		bot,
		logger,
		services,
	}
}

func (b *Bot) setListBotCommands(ctx context.Context) {
	commands := &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "/start", Description: "Команда для запуска бота"},
		},
	}
	_, err := b.Instance.SetMyCommands(ctx, commands)
	if err != nil {
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}
}

func (b *Bot) SetHandlers(ctx context.Context) {
	startCmd := commands.Start(b.logger, b.services.Profile)
	b.Instance.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startCmd)

	addChartCb := scripts.CreateChart(b.logger, b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandler(bot.HandlerTypeCallbackQueryData, "createchart", bot.MatchTypeExact, addChartCb)

	getUserChartsCb := callbacks.UserCharts(b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandler(bot.HandlerTypeCallbackQueryData, "charts", bot.MatchTypeExact, getUserChartsCb)

	regexpAdvancedChrt := regexp.MustCompile(`AdvancedChrt: [0-9]+:deleteTo:[0-9]+`)
	advancedChrtCb := callbacks.AdvancedChrt(b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandlerRegexp(bot.HandlerTypeCallbackQueryData, regexpAdvancedChrt, advancedChrtCb)

	regexpCreateForecast := regexp.MustCompile(`CreateForecast:[0-9]+`)
	createForecastCb := scripts.CreateForecast(b.logger, b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandlerRegexp(bot.HandlerTypeCallbackQueryData, regexpCreateForecast, createForecastCb)

	regexpInterpritationChart := regexp.MustCompile(`InterpritationChart: [0-9]+`)
	interpritaionChartCb := callbacks.InterpritaionChrt(b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandlerRegexp(bot.HandlerTypeCallbackQueryData, regexpInterpritationChart, interpritaionChartCb)

	regexpSelectChartParam := regexp.MustCompile(`SelectChartParam:[0-9]+:[0-9]+`)
	interpritaionSelectChartParamCb := callbacks.SelectChartParam(b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandlerRegexp(bot.HandlerTypeCallbackQueryData, regexpSelectChartParam, interpritaionSelectChartParamCb)

	b.setListBotCommands(ctx)
}

func (b *Bot) StartViaWebhook(ctx context.Context, url string) {
	op := "internal.bot.Bot.StartViaWebHook"
	logger := b.logger.With(
		slog.String("func", op),
		slog.String("HookURL", url),
	)

	_, err := b.Instance.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: url,
	})
	if err != nil {
		logger.Error("Error while setting the webhok for the bot", err)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		return
	}

	go b.Instance.StartWebhook(ctx)

	logger.Info("Bot has started")
}

func (b *Bot) StopWithWebHook(ctx context.Context) {
	op := "internal.bot.Bot.StopWithWebHook"
	logger := b.logger.With(
		slog.String("func", op),
	)

	_, err := b.Instance.DeleteWebhook(ctx, &bot.DeleteWebhookParams{
		DropPendingUpdates: true,
	})
	if err != nil {
		logger.Error("Occured the error while stopping the bot", err)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		return
	}

	logger.Info("Bot was successfully stopped")
}
