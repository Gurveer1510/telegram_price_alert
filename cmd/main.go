package main

import (
	"fmt"
	"time"

	"github.com/Gurveer1510/telegram_price_tracker/internal/config"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

var (
	ticker *kiteticker.Ticker
)

var (
	instToken = []uint32{138095876, 539437}
)

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	fmt.Println("Connected")
	err := ticker.Subscribe(instToken)
	if err != nil {
		fmt.Println("err: ", err)
	}
	// Set subscription mode for the subscribed token
	// Default mode is Quote
	err = ticker.SetMode(kiteticker.ModeFull, instToken)
	if err != nil {
		fmt.Println("err: ", err)
	}

}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {
	fmt.Printf("Tick: %+v", tick)
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	fmt.Printf("Order: %s", order.OrderID)
}

func main() {
	cfg, _ := config.GetConfig()

	// dsn := db.DSN(cfg)
	// db, err := db.NewPool(context.Background(), dsn)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer db.Pool.Close()

	// repo := repository.NewTelegramRepo(db)
	// usecase := usecase.NewTelegramUseCase(repo)

	// tgBot, err := telegram.Newbot(cfg.BotToken, usecase)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// for {
	// 	tgBot.GetUpdates()
	// 	log.Println("Bot disconnected, restarting...")
	// 	time.Sleep(5 * time.Second)
	// }

	apiKey := cfg.ZerodhaApiKey
	accessToken := "q1F8jDeeOP9qTlB6O75XGosHpzYShJ3n"
	// zerodhaclient := zerodha.NewZerodhaClient(accessToken, apiKey)
	// zerodhaclient.
	// zerodhaclient.GetInstruments()
	// Create new Kite ticker instance
	ticker := kiteticker.New(apiKey, accessToken)

	// // Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// // Start the connection
	ticker.Serve()

}
