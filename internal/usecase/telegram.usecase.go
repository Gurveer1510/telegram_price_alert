package usecase

import (
	"context"
	"log"
	"strconv"

	"github.com/Gurveer1510/telegram_price_tracker/internal/repository"
	"github.com/Gurveer1510/telegram_price_tracker/internal/types"
	"github.com/Gurveer1510/telegram_price_tracker/internal/utils"
	"github.com/Gurveer1510/telegram_price_tracker/internal/zerodha"
)

type TelegramUsecase struct {
	repo          *repository.TelegramZerodhaRepo
	tickerService *zerodha.ZerodhaTicker
	zerodhaClient *zerodha.ZerodhaClient
}

func NewTelegramUseCase(repo *repository.TelegramZerodhaRepo, tickerSvc *zerodha.ZerodhaTicker) *TelegramUsecase {
	return &TelegramUsecase{repo: repo, tickerService: tickerSvc}
}

func (tuc *TelegramUsecase) CreateAlert(ctx context.Context, alert *types.Alert) (int64, error) {
	token := utils.GetToken(alert.Instrument_name)
	if token == "" {
		log.Println("Token not found for instrument ", alert.Instrument_name)
		return 0, nil
	}
	log.Println("Token found for instrument ", alert.Instrument_name, " is ", token)

	instrumentToken, err := strconv.Atoi(token)
	if err != nil {
		log.Println(err.Error())
	}
	alert.Instrument_token = instrumentToken
	alert.Exchange = "NSE"
	id, err := tuc.repo.Create(ctx, alert)
	if err != nil {
		return 0, nil
	}
	tuc.tickerService.Subscribe(uint32(alert.Instrument_token))
	return id, nil
}

func (tuc *TelegramUsecase) StoreAccessToken(ctx context.Context, token string) error {
	return tuc.repo.StoreAccessToken(ctx, token)
}

func (tuc *TelegramUsecase) GetLatestAccessToken(ctx context.Context) (string, error) {
	return tuc.repo.GetLatestAccessToken(ctx)
}
