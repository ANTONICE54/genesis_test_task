package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GetRateResponse struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

// Function for fetching exchange rate from third-party API
func FetchRateData(apiKey string) (*float64, error) {

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", apiKey)
	rateInfo, err := http.Get(url)
	if err != nil {

		return nil, err
	}

	body, err := io.ReadAll(rateInfo.Body)
	if err != nil {

		return nil, err
	}

	var rateList GetRateResponse
	err = json.Unmarshal(body, &rateList)
	if err != nil {
		return nil, err
	}

	rate := rateList.ConversionRates["UAH"]

	return &rate, nil
}

// The Server.cronOperator accepts this string as a parameter to set the time at which emails will be sent
func TimeToSendEmails(time string) string {

	splittedTime := strings.Split(time, ":")

	return fmt.Sprintf("%s %s %s * * *", splittedTime[2], splittedTime[1], splittedTime[0])
}
