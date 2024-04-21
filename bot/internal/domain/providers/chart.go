package providers

import (
	"context"

	"github.com/Izumra/Magistus/bot/internal/domain/dto"
	"github.com/Izumra/Magistus/bot/internal/domain/entity"
)

type Chart interface {
	ChartById(ctx context.Context, IdChart int64) (*entity.Chart, error)
	ListChartsByIdUser(ctx context.Context, id_user int64) ([]*dto.SimpChart, error)
}
