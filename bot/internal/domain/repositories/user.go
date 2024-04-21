package repositories

import (
	"context"

	"github.com/go-telegram/bot/models"
)

type User interface {
	AddUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id_user int64) error
}
