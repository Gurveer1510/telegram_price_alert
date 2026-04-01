package repository

import (
	"context"
	"time"

	"github.com/Gurveer1510/telegram_price_tracker/internal/db"
	"github.com/Gurveer1510/telegram_price_tracker/internal/types"
)

type TelegramZerodhaRepo struct {
	db *db.DB
}

func NewTelegramZerodhaRepo(db *db.DB) *TelegramZerodhaRepo {
	return &TelegramZerodhaRepo{db: db}
}

func (tr *TelegramZerodhaRepo) Create(ctx context.Context, alert *types.Alert) (int64, error) {
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

func (tr *TelegramZerodhaRepo) StoreAccessToken(ctx context.Context, accessToken string) error {
	query := `INSERT INTO access_tokens(access_token, created_at) VALUES($1, $2)`
	createdAt := time.Now()
	_, err := tr.db.Pool.Exec(ctx, query, accessToken, createdAt)
	return err
}

func (tr *TelegramZerodhaRepo) GetLatestAccessToken(ctx context.Context) (string, error) {
	query := `	SELECT * FROM access_tokens ORDER BY created_at DESC LIMIT 1 RETURNING access_token;`
	var token string
	err := tr.db.Pool.QueryRow(ctx, query).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
