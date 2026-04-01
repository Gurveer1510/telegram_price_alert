package zerodha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ZerodhaClient struct {
	ApiKey      string
	AccessToken string
}

func NewZerodhaClient(apiKey, accessToken string) *ZerodhaClient {
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

type LoginPayload struct {
	Username	string
	Password	string
}

func (z *ZerodhaClient) LoginWithUserPassZerodha(username, password string) error {

	urlString := "https://kite.zerodha.com/api/login"

	body := LoginPayload{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(body)

	resp, err := http.Post(urlString, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	log.Println("RESPONSE BODY \n ", string(bodyBytes))

	return nil
}