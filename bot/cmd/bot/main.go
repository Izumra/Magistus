package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

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

	bot.Start(ctx)

	chanSignals := make(chan os.Signal, 1)
	signal.Notify(chanSignals, os.Interrupt)

	select {
	case sign := <-chanSignals:
		bot.Stop(ctx)
		logger.Info("Program was finished", slog.Any("sign", sign))
	}
}
