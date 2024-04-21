package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Izumra/Magistus/bot/internal/domain/dto"
	"github.com/Izumra/Magistus/bot/internal/domain/dto/sotis"
	"github.com/Izumra/Magistus/bot/internal/domain/entity"
	"github.com/Izumra/Magistus/bot/internal/storage"
)

func (s *Storage) ChartById(ctx context.Context, IdChart int64) (*entity.Chart, error) {
	op := "storage/sqlite/ChartStorage.ChartById"

	query := "select * from charts where id=?"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results, err := state.QueryContext(ctx, IdChart)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrChartNotFound
		}
		return nil, err
	}

	var chart entity.Chart
	for results.Next() {
		if results.Err() != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		results.Scan(&chart.Id, &chart.Query, &chart.Content, &chart.IdCreator, &chart.Title, &chart.Interpritaion)
	}

	return &chart, nil
}

func (s *Storage) ListChartsByIdUser(ctx context.Context, id_user int64) ([]*dto.SimpChart, error) {
	op := "storage/sqlite/ChartStorage.ListChartsByIdUser"

	query := "select id,title from charts where id_creator=?"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results, err := state.QueryContext(ctx, id_user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserChartsNotFound
		}
		return nil, err
	}

	var charts []*dto.SimpChart
	for results.Next() {
		if results.Err() != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		var chart dto.SimpChart
		results.Scan(&chart.Id, &chart.Title)
		charts = append(charts, &chart)
	}

	return charts, nil
}

func (s *Storage) AddChart(ctx context.Context, nameChart string, chartData *sotis.CreateChartResp, IdCreator int64) (int64, error) {
	op := "storage/sqlite/ChartStorage.AddChart"

	query := "insert into charts(title,query,content,id_creator,interpritation)values(?,?,?,?,?)"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := state.ExecContext(ctx, nameChart, chartData.Query, chartData.Content, IdCreator, "")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateInterpritation(ctx context.Context, id_chart int64, interpritation *string) error {
	op := "storage/sqlite/ChartStorage.UpdateInterpritation"

	query := "update charts set interpritation=? where id=?"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = state.ExecContext(ctx, interpritation, id_chart)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
