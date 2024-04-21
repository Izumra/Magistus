package repositories

import (
	"context"

	"github.com/Izumra/Magistus/bot/internal/domain/dto/sotis"
)

type Chart interface {
	AddChart(ctx context.Context, nameChart string, chartData *sotis.CreateChartResp, IdCreator int64) (int64, error)
	UpdateInterpritation(ctx context.Context, id_chart int64, interpritation *string) error
}
