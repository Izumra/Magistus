package chart

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/Izumra/Magistus/bot/internal/domain/dto/sotis"
	"github.com/Izumra/Magistus/bot/internal/domain/entity"
	"github.com/Izumra/Magistus/bot/internal/domain/providers"
	"github.com/Izumra/Magistus/bot/internal/domain/repositories"
	"github.com/Izumra/Magistus/bot/internal/lib/converter"
	"github.com/Izumra/Magistus/bot/internal/lib/req"
	"github.com/Izumra/Magistus/bot/internal/storage"
)

var (
	ErrCreateChart    = errors.New("✖️ Не удалось построить натальную карту")
	ErrCreateForecast = errors.New("✖️ Не удалось составить прогноз")
	ErrChartNotFound  = errors.New("✖️ Натальная карта не найдена")
	ErrUnexpected     = errors.New("✖️ Неизвестная ошибка, попробуйте еще раз")
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

func (s *Service) GetChart(ctx context.Context, id_chart int64) (*entity.Chart, error) {
	op := "services.chart.Service.GetChart"
	logger := s.log.With(slog.String("method", op))

	chart, err := s.chartPrvdr.ChartById(ctx, id_chart)
	if err != nil {
		if errors.Is(err, storage.ErrChartNotFound) {
			return nil, ErrChartNotFound
		}
		logger.Error("occured the error while getting the chart", err)
		return nil, ErrUnexpected
	}

	return chart, nil
}

func (s *Service) CreateChart(
	ctx context.Context,
	IdCreator int64,
	name string,
	dateBorn time.Time,
	cords *converter.MapCords,
) (int64, error) {
	op := "services.chart.Service.CreateChart"
	logger := s.log.With(slog.String("method", op))

	buffer := &bytes.Buffer{}
	form := multipart.NewWriter(buffer)

	strDateBorn := dateBorn.String()

	form.WriteField("ct0", "0")
	form.WriteField("name0", name)
	form.WriteField("year0", strDateBorn[0:4])
	form.WriteField("month0", strDateBorn[5:7])
	form.WriteField("day0", strDateBorn[8:10])
	form.WriteField("hour0", strDateBorn[11:13])
	form.WriteField("min0", strDateBorn[14:16])
	form.WriteField("sec0", strDateBorn[17:19])
	form.WriteField("gmt0", string(strDateBorn[20:25]))
	form.WriteField("lat_d0", fmt.Sprint(cords.LatDeg))
	form.WriteField("lat_m0", fmt.Sprint(cords.LatMin))
	form.WriteField("lat_s0", fmt.Sprint(cords.LatSec))
	form.WriteField("lat_0", fmt.Sprint(cords.LatDir))
	form.WriteField("lon_d0", fmt.Sprint(cords.LongDeg))
	form.WriteField("lon_m0", fmt.Sprint(cords.LongMin))
	form.WriteField("lon_s0", fmt.Sprint(cords.LongSec))
	form.WriteField("lon_0", fmt.Sprint(cords.LongDir))

	form.Close()

	headers := map[string]string{
		"Content-Type": form.FormDataContentType(),
	}

	var resp sotis.CreateChartResp

	err := req.WithMultipartFormData(
		ctx,
		"POST",
		fmt.Sprintf("https://sotis-online.ru/edit1.php/?chr=dt:%v%v%v%v%v%v;%s",
			strDateBorn[0:4],
			strDateBorn[5:7],
			strDateBorn[8:10],
			strDateBorn[11:13],
			strDateBorn[14:16],
			strDateBorn[17:19],
			cords.CordsReqString,
		),
		headers,
		buffer,
		&resp,
	)

	if err != nil {
		logger.Error("Occured the error while building the chart", slog.Any("error", err))
		return -1, ErrCreateChart
	}
	if resp.Content == "" || resp.Query == "" {
		return -1, ErrCreateChart
	}

	resp.Query += fmt.Sprintf(";%s", cords.CordsReqString)
	IdChart, err := s.chartRep.AddChart(ctx, name, &resp, IdCreator)
	if err != nil {
		logger.Error("Occured the error while additing the chart to db", slog.Any("error", err))
		return -1, ErrCreateChart
	}

	go s.updateChrtInterpritaion(ctx, IdChart)

	return IdChart, nil
}

func (s *Service) CreateForecast(
	ctx context.Context,
	id_chart int64,
	date_forecast time.Time,
	cords_forecast *converter.MapCords,
) (string, error) {
	op := "services.chart.Service.CreateForecast"
	logger := s.log.With(slog.String("method", op))

	chart, err := s.GetChart(ctx, id_chart)
	if err != nil {
		return "", err
	}
	if len(chart.Query) < 6 {
		return "", ErrCreateForecast
	}

	strDateForecast := date_forecast.String()
	gmt := string([]byte(strDateForecast)[20:25])

	buffer := &bytes.Buffer{}
	form := multipart.NewWriter(buffer)

	interprit := string([]byte(chart.Query)[5:]) + "|ct:1;" + fmt.Sprintf("dt:%v%v%v%v%v%v;%s;gmt:%s",
		strDateForecast[0:4],
		strDateForecast[5:7],
		strDateForecast[8:10],
		strDateForecast[11:13],
		strDateForecast[14:16],
		strDateForecast[17:19],
		cords_forecast.CordsReqString,
		gmt,
	)

	log.Println(interprit)

	form.WriteField("chr", interprit)

	form.Close()

	headers := map[string]string{
		"Content-Type": form.FormDataContentType(),
	}

	var resp string

	err = req.WithMultipartFormData(
		ctx,
		"POST",
		"https://sotis-online.ru/addTract.php",
		headers,
		buffer,
		&resp,
	)
	if err != nil {
		logger.Error("Occured the error while interpritating the chart", slog.Any("error", err))
		return "", err
	}

	return resp, nil
}

func (s *Service) updateChrtInterpritaion(ctx context.Context, IdChart int64) {
	op := "services.chart.Service.CreateChart"
	logger := s.log.With(slog.String("method", op))

	chart, err := s.chartPrvdr.ChartById(ctx, IdChart)
	if err != nil {
		logger.Info("Chart was't found by cause", err)
		return
	}

	buffer := &bytes.Buffer{}
	form := multipart.NewWriter(buffer)

	form.WriteField("chr", string([]byte(chart.Query)[5:]))

	form.Close()

	headers := map[string]string{
		"Content-Type": form.FormDataContentType(),
	}

	var resp string

	err = req.WithMultipartFormData(
		ctx,
		"POST",
		"https://sotis-online.ru/addTract.php",
		headers,
		buffer,
		&resp,
	)
	if err != nil {
		logger.Error("Occured the error while interpritating the chart", slog.Any("error", err))
		return
	}

	err = s.chartRep.UpdateInterpritation(ctx, IdChart, &resp)
	if err != nil {
		logger.Error("Occured the error while generate the chart interpritation", slog.Any("error", err))
	}
}
