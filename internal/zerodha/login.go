package zerodha

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Gurveer1510/telegram_price_tracker/internal/config"
)

func LoginWithZerodha() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	resp, err := http.Get(fmt.Sprintf("https://kite.zerodha.com/connect/login?v=3&api_key=%v", cfg.ZerodhaApiKey))
	if err != nil {
		return err
	}

	for key, val := range resp.Header {
		log.Println(key, val)
	}
	fmt.Println("---------------------------------------------------")
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(string(bodyBytes))
	fmt.Println("----------------------------------------------------")
	log.Println(resp.StatusCode)
	defer resp.Body.Close()
	return nil
}
