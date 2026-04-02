package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/Gurveer1510/telegram_price_tracker/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
)

type AlertChecker struct {
	repo *repository.TelegramZerodhaRepo
	bot  *tgbotapi.BotAPI
}

func NewAlertChecker(repo *repository.TelegramZerodhaRepo, bot *tgbotapi.BotAPI) *AlertChecker {
	return &AlertChecker{repo: repo, bot: bot}
}

func (ac *AlertChecker) CheckAlerts(tick kitemodels.Tick) {
	ctx := context.Background()
	price := tick.LastPrice

	alerts, err := ac.repo.GetAlerts(ctx, tick.InstrumentToken)
	if err != nil {
		log.Println("CheckAlerts query error:", err)
		return
	}

	for _, alert := range alerts {
		triggered := (alert.Condition == "above" && price > alert.Trigger_price) || (alert.Condition == "below" && price < alert.Trigger_price)
		if !triggered {
			continue
		}

		msg := tgbotapi.NewMessage(alert.ChatId, fmt.Sprintf("🔔 %v hit %.2f (your alert: %v %.2f)", alert.Instrument_name, price, alert.Condition, alert.Trigger_price))
		ac.bot.Send(msg)
		ac.repo.DeleteAlert(ctx, alert.ID)
	}

}
