package main

import (
	"context"
	"log"
	"time"

	"github.com/Gurveer1510/telegram_price_tracker/internal/config"
	"github.com/Gurveer1510/telegram_price_tracker/internal/db"
	"github.com/Gurveer1510/telegram_price_tracker/internal/repository"
	"github.com/Gurveer1510/telegram_price_tracker/internal/telegram"
	"github.com/Gurveer1510/telegram_price_tracker/internal/usecase"
	"github.com/Gurveer1510/telegram_price_tracker/internal/zerodha"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	database, err := db.NewPool(context.Background(), db.DSN(cfg))
	if err != nil {
		log.Println(err)
		return
	}

	defer database.Pool.Close()

	zClient := zerodha.NewZerodhaClient(cfg.ZerodhaApiKey, cfg.KiteUser, cfg.KitePassword, cfg.TotpSecret, cfg.KiteSecret)
	repo := repository.NewTelegramZerodhaRepo(database)

	tgBot, err := telegram.Newbot(cfg.BotToken, nil)
	if err != nil {
		log.Println(err)
		return
	}

	alertCheck := usecase.NewAlertChecker(repo, tgBot.Bot, nil)
	ticker := zerodha.NewZerodhaTicker(cfg.ZerodhaApiKey, zClient.AccessToken, alertCheck)
	existingAlerts, _ := repo.GetAllAlerts(context.TODO())
	for _, token := range existingAlerts {
		ticker.Subscribe(uint32(token.Instrument_token))
	}
	alertCheck.ZerodhaTicker = ticker
	uc := usecase.NewTelegramUseCase(repo, ticker)
	tgBot.TGUsecase = uc
	ticker.Start()

	if err := uc.StoreAccessToken(context.Background(), zClient.AccessToken); err != nil {
		log.Println(err)
	}
	for {
		tgBot.GetUpdates()
		log.Println("Bot disconnected, restarting...")
		time.Sleep(5 * time.Second)
	}

}
