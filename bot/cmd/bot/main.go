package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Izumra/Magistus/bot/internal/app"
	"github.com/Izumra/Magistus/bot/internal/config"
	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
	"github.com/Izumra/Magistus/bot/internal/storage/sqlite"
	"github.com/Izumra/Magistus/bot/lib/logger"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	logger := logger.New(logger.Local, os.Stdout)

	db := sqlite.New("sqlite3", cfg.Db.Path)
	logger.Info("Programm has successfully connected to the db")

	services := &app.Services{
		Chart:   chart.New(logger, db, db, db, db),
		Profile: profile.New(logger, db, db, db, db),
	}

	bot := app.New(logger, cfg.Bot.Token, services)

	bot.SetHandlers(ctx)

	bot.StartViaWebhook(ctx, cfg.Bot.Webhook)

	go http.ListenAndServe(":3222", bot.Instance.WebhookHandler())

	chanSignals := make(chan os.Signal, 1)
	signal.Notify(chanSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sign := <-chanSignals
	logger.Info("Program was finished with signal", sign)
}
