package zerodha

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type ZerodhaClient struct {
	ApiKey      string
	AccessToken string
}

func NewZerodhaClient(accessToken, apiKey string) *ZerodhaClient {
	return &ZerodhaClient{
		ApiKey:      apiKey,
		AccessToken: accessToken,
	}
}

func (z *ZerodhaClient) LoginWithZerodha() error {

	resp, err := http.Get(fmt.Sprintf("https://kite.zerodha.com/connect/login?v=3&api_key=%v", z.ApiKey))
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
