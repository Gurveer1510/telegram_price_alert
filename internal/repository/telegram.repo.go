package repository

import (
	"context"

	"github.com/Gurveer1510/telegram_price_tracker/internal/db"
	"github.com/Gurveer1510/telegram_price_tracker/internal/types"
)

type TelegramRepo struct {
	db *db.DB
}

func NewTelegramRepo(db *db.DB) *TelegramRepo {
	return &TelegramRepo{db: db}
}

func (tr *TelegramRepo) Create(ctx context.Context, alert *types.Alert) (int64, error) {
	query := `
		INSERT INTO alerts(instrument_token, instrument_name, chat_id, exchange, trigger_price, condition) VALUES($1, $2, $3, $4, $5, $6) RETURNING id
	`
	var id int64
	err := tr.db.Pool.QueryRow(ctx, query, alert.Instrument_token, alert.Instrument_name, alert.ChatId, alert.Exchange, alert.Trigger_price, alert.Condition).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}
