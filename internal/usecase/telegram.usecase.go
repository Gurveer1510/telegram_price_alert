package usecase

import (
	"context"

	"github.com/Gurveer1510/telegram_price_tracker/internal/repository"
	"github.com/Gurveer1510/telegram_price_tracker/internal/types"
)

type TelegramUsecase struct {
	repo repository.TelegramRepo
}

func NewTelegramUseCase(repo *repository.TelegramRepo) *TelegramUsecase {
	return &TelegramUsecase{repo: *repo}
}

func (tuc *TelegramUsecase) CreateAlert(ctx context.Context, alert *types.Alert) (int64, error) {
	return tuc.repo.Create(ctx, alert)
}
