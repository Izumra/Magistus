package profile

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Izumra/Magistus/bot/internal/domain/dto"
	"github.com/Izumra/Magistus/bot/internal/domain/providers"
	"github.com/Izumra/Magistus/bot/internal/domain/repositories"
	"github.com/Izumra/Magistus/bot/internal/storage"
	"github.com/go-telegram/bot/models"
)

var (
	ErrUnexpected     = errors.New("✖️ Неизвестная ошибка, попробуйте еще раз")
	ErrChartsNotFound = errors.New("✖️ Ни одной карты еще не создано")
)

type Service struct {
	log        *slog.Logger
	usrPrvdr   providers.User
	chartPrvdr providers.Chart
	usrRep     repositories.User
	chartRep   repositories.Chart
}

func New(
	log *slog.Logger,
	usrPrvdr providers.User,
	chartPrvdr providers.Chart,
	usrRep repositories.User,
	chartRep repositories.Chart,
) *Service {
	return &Service{
		log,
		usrPrvdr,
		chartPrvdr,
		usrRep,
		chartRep,
	}
}

func (s *Service) AuthProfile(
	ctx context.Context,
	user *models.User,
) error {
	op := "services.profile.Service.AuthProfile"
	logger := s.log.With(slog.String("method", op))

	_, err := s.usrPrvdr.UserByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			err := s.usrRep.AddUser(ctx, user)
			if err != nil {
				logger.Error("Error while additing the user", err)
				return err
			}

			return nil
		}
		logger.Error("Error while finding the user", err)

		return err
	}

	return nil
}

func (s *Service) DeleteProfile(ctx context.Context, IdProfile int64) error {
	return s.usrRep.DeleteUser(ctx, IdProfile)
}

func (s *Service) ListCharts(
	ctx context.Context,
	id_user int64,
) ([]*dto.SimpChart, error) {
	op := "services.profile.Service.ListCharts"
	logger := s.log.With(slog.String("method", op))

	charts, err := s.chartPrvdr.ListChartsByIdUser(ctx, id_user)
	if err != nil {
		if errors.Is(err, storage.ErrChartsNotFound) {
			return nil, ErrChartsNotFound
		}
		logger.Error("occured the error while getting profile charts", err)
		return nil, ErrUnexpected
	}

	return charts, nil
}
