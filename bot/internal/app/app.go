package app

import (
	"context"
	"log/slog"
	"regexp"
	"syscall"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Izumra/Magistus/bot/internal/app/handlers/callbacks"
	"github.com/Izumra/Magistus/bot/internal/app/handlers/commands"
	"github.com/Izumra/Magistus/bot/internal/app/handlers/scripts"
	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
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
	b.Instance.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"createchart",
		bot.MatchTypeExact,
		addChartCb,
	)

	getUserChartsCb := callbacks.UserCharts(b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"charts",
		bot.MatchTypeExact,
		getUserChartsCb,
	)

	regexpAdvancedChrt := regexp.MustCompile(`AdvancedChrt:[0-9]+(:deleteTo:[0-9]+)?`)
	advancedChrtCb := callbacks.AdvancedChrt(b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandlerRegexp(
		bot.HandlerTypeCallbackQueryData,
		regexpAdvancedChrt,
		advancedChrtCb,
	)

	regexpCreateForecast := regexp.MustCompile(`CreateForecast:[0-9]+`)
	createForecastCb := scripts.CreateForecast(b.logger, b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandlerRegexp(
		bot.HandlerTypeCallbackQueryData,
		regexpCreateForecast,
		createForecastCb,
	)

	regexpRemoveChart := regexp.MustCompile(`RemoveChart:[0-9]+`)
	removeChartCb := callbacks.RemoveChart(b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandlerRegexp(
		bot.HandlerTypeCallbackQueryData,
		regexpRemoveChart,
		removeChartCb,
	)

	regexpInterpritationChart := regexp.MustCompile(`InterpritationChart:[0-9]+`)
	interpritaionChartCb := callbacks.InterpritaionChrt(b.services.Profile, b.services.Chart)
	b.Instance.RegisterHandlerRegexp(
		bot.HandlerTypeCallbackQueryData,
		regexpInterpritationChart,
		interpritaionChartCb,
	)

	regexpSelectChartParam := regexp.MustCompile(`SelectChartParam:[0-9]+:[0-9]+`)
	interpritaionSelectChartParamCb := callbacks.SelectChartParam(
		b.services.Profile,
		b.services.Chart,
	)
	b.Instance.RegisterHandlerRegexp(
		bot.HandlerTypeCallbackQueryData,
		regexpSelectChartParam,
		interpritaionSelectChartParamCb,
	)

	b.setListBotCommands(ctx)
}

func (b *Bot) Start(ctx context.Context) {
	op := "internal.bot.Bot.Start"
	logger := b.logger.With(
		slog.String("func", op),
	)

	go b.Instance.Start(ctx)

	logger.Info("Bot has started", slog.Int("pid", syscall.Getpid()))
}

func (b *Bot) Stop(ctx context.Context) {
	op := "internal.bot.Bot.Stop"
	logger := b.logger.With(
		slog.String("func", op),
	)

	_, err := b.Instance.DeleteWebhook(ctx, &bot.DeleteWebhookParams{
		DropPendingUpdates: true,
	})
	if err != nil {
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		return
	}

	logger.Info("Bot was successfully stopped")
}
