package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func GetToken(symbol string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	file, err := os.Open(cwd + "/internal/instrument_dump/zerodha.csv")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			log.Println(err.Error())
			break
		}
		if err != nil {
			log.Println("Error reading now:", err)
		}
		if record[2] == strings.ToUpper(symbol) {
			log.Println("Found the token for ", symbol, " and the token is ", record[0])
			return record[0]
		}
	}

	log.Println("Token not found for symbol ", symbol)
	return ""
}
