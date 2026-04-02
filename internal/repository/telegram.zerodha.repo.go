package repository

import (
	"context"
	"log"
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

// TODO: Implement GetAlerts to fetch all alerts for a given instrument token and exchange
func (tr *TelegramZerodhaRepo) GetAlerts(ctx context.Context, token uint32) ([]types.Alert, error) {
	query := `SELECT id, instrument_token, instrument_name, chat_id, exchange, trigger_price, condition 
		FROM alerts WHERE instrument_token=$1`

	rows, err := tr.db.Pool.Query(ctx, query, token)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var alerts []types.Alert
	for rows.Next() {
		var a types.Alert
		err := rows.Scan(&a.ID, &a.Instrument_token, &a.Instrument_name, &a.ChatId, &a.Exchange, &a.Trigger_price, &a.Condition)
		if err != nil {
			log.Println("ERROR WHILE SCANNING ROWS",err.Error())
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (tr *TelegramZerodhaRepo) GetAllAlerts(ctx context.Context) ([]types.Alert, error) {
	query := `SELECT id, instrument_token, instrument_name, chat_id, exchange, trigger_price, condition 
		FROM alerts `

	rows, err := tr.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var alerts []types.Alert
	for rows.Next() {
		var a types.Alert
		err := rows.Scan(&a.ID, &a.Instrument_token, &a.Instrument_name, &a.ChatId, &a.Exchange, &a.Trigger_price, &a.Condition)
		if err != nil {
			log.Println("ERROR WHILE SCANNING ROWS",err.Error())
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (tr *TelegramZerodhaRepo) DeleteAlert(ctx context.Context, id int) error {
	query := `DELETE FROM alerts WHERE id = $1`
	_, err := tr.db.Pool.Exec(ctx, query, id)
	return err
}

func (tr *TelegramZerodhaRepo) StoreAccessToken(ctx context.Context, accessToken string) error {
	query := ` INSERT INTO access_tokens(id, access_token, created_at)
        VALUES(1, $1, $2)
        ON CONFLICT (id) DO UPDATE
            SET access_token = EXCLUDED.access_token,
                created_at   = EXCLUDED.created_at`
	createdAt := time.Now()
	_, err := tr.db.Pool.Exec(ctx, query, accessToken, createdAt)
	return err
}

func (tr *TelegramZerodhaRepo) GetLatestAccessToken(ctx context.Context) (string, error) {
    var token string
    err := tr.db.Pool.QueryRow(ctx, `SELECT access_token FROM access_tokens WHERE id = 1`).Scan(&token)
    return token, err
}
