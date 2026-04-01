package usecase

import (
	"context"

	"github.com/Gurveer1510/telegram_price_tracker/internal/repository"
	"github.com/Gurveer1510/telegram_price_tracker/internal/types"
)

type TelegramUsecase struct {
	repo repository.TelegramZerodhaRepo
}

func NewTelegramUseCase(repo *repository.TelegramZerodhaRepo) *TelegramUsecase {
	return &TelegramUsecase{repo: *repo}
}

func (tuc *TelegramUsecase) CreateAlert(ctx context.Context, alert *types.Alert) (int64, error) {
	return tuc.repo.Create(ctx, alert)
}

func (tuc *TelegramUsecase) StoreAccessToken(ctx context.Context, token string) (error) {
	return tuc.repo.StoreAccessToken(ctx, token)
}

func (tuc *TelegramUsecase) GetLatestAccessToken(ctx context.Context) (string, error) {
	return tuc.repo.GetLatestAccessToken(ctx)
}