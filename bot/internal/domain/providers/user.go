package providers

import (
	"context"

	"github.com/Izumra/Magistus/bot/internal/domain/entity"
)

type User interface {
	UserByID(ctx context.Context, id int64) (*entity.User, error)
	AllUsers(ctx context.Context) ([]*entity.User, error)
}
