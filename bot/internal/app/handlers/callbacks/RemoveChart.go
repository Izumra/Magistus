package callbacks

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Izumra/Magistus/bot/internal/services/chart"
	"github.com/Izumra/Magistus/bot/internal/services/profile"
)

func RemoveChart(
	prof *profile.Service,
	chart *chart.Service,
) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chartActionsData := strings.Split(update.CallbackQuery.Data, ":")
		chartId, err := strconv.ParseInt(chartActionsData[1], 10, 0)
		if err == nil {
			err = chart.DeleteChart(ctx, chartId)
			if err != nil {
				log.Println(err)
				return
			}

			UserChartsHandler(prof, chart, ctx, b, update)

		}
	}
}
