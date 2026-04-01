package zerodha

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (z *ZerodhaClient) GetInstruments() error {
	url := "https://api.kite.trade/instruments"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	authHeader := fmt.Sprintf("token %v:%v", z.ApiKey,z.AccessToken)
	// Headers
	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	cwd, _ := os.Getwd()
	path := cwd+"/internal/instrument_dump/"
	file, err := os.Create(path+"zerodha.csv")
	if err != nil {
		return err
	}

	n, err := file.Write(body)
	if err != nil {
		return err
	}
	log.Println(n," instruments found.")
	return nil
}
