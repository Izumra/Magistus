package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Izumra/Magistus/bot/internal/domain/entity"
	"github.com/Izumra/Magistus/bot/internal/storage"
	"github.com/go-telegram/bot/models"
)

func (s *Storage) UserByID(ctx context.Context, id int64) (*entity.User, error) {
	op := "storage/sqlite/UserStorage.UserByID"

	query := "select * from users where id=?"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results, err := state.QueryContext(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}

	var user entity.User
	for results.Next() {
		if results.Err() != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		results.Scan(&user.Id, &user.Username)
	}

	return &user, nil
}

func (s *Storage) AllUsers(ctx context.Context) ([]*entity.User, error) {
	op := "storage/sqlite/UserStorage.AllUsers"

	query := "select * from users"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results, err := state.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUsersNotFound
		}
		return nil, err
	}

	var users []*entity.User
	for results.Next() {
		if results.Err() != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		var user entity.User
		results.Scan(&user.Id, &user.Username)
		users = append(users, &user)
	}

	return users, nil
}

func (s *Storage) AddUser(ctx context.Context, user *models.User) error {
	op := "storage/sqlite/UserStorage.AddUser"

	query := "insert into users(id,username)values(?,?)"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = state.ExecContext(ctx, user.ID, user.Username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id_user int64) error {
	op := "storage/sqlite/UserStorage.DeleteUser"

	query := "delete from users where id=?"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = state.ExecContext(ctx, id_user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
