package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CoinGeckoData struct {
	ID         string `json:"id"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	MarketData struct {
		CurrentPrice struct {
			Eur float64 `json:"eur"`
			Gbp float64 `json:"gbp"`
			Inr float64 `json:"inr"`
			Usd float64 `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}

const (
	PriceURL     = "https://api.coingecko.com/api/v3/coins/"
	DefaultCrypt = "bitcoin"
)

func main() {
	crypto := promptCrypto()
	if crypto == "" {
		fmt.Println("You have not given an input. Checking the current value of Bitcoin.")
		fmt.Println()
		crypto = DefaultCrypt
	} else {
		fmt.Println("Please wait....")
		fmt.Println()
	}
	priceCheck(crypto)
}

func promptCrypto() string {
	fmt.Println("Which crypto price do you want to check? eg: ethereum, bitcoin, tether, usd-coin...etc")
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.ToLower(strings.TrimSpace(scanner.Text()))
}

func priceCheck(c string) {
	crypto := strings.ReplaceAll(strings.Trim(c, " "), " ", "-")
	coin, err := fetchCoinData(crypto)
	if err != nil {
		fmt.Printf("Crypto not found! Please check if \"%s\" is a valid crypto.\n", cases.Title(language.English).String(c))
		return
	}
	printCoinPrice(coin)
}

func fetchCoinData(cryptoID string) (*CoinGeckoData, error) {
	resp, err := http.Get(PriceURL + cryptoID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var coin CoinGeckoData
	if err := json.Unmarshal(body, &coin); err != nil {
		return nil, err
	}
	return &coin, nil
}

func printCoinPrice(coin *CoinGeckoData) {
	now := time.Now()
	fmt.Printf("%s (%s) Price Today (%v %v, %v | %v):\n",
		coin.Name, strings.ToUpper(coin.Symbol),
		now.Month(), now.Day(), now.Year(), now.Format(time.Kitchen),
	)
	fmt.Printf("USD: %.2f\nEUR: %.2f\nGBP: %.2f\nINR: %.2f\n",
		coin.MarketData.CurrentPrice.Usd,
		coin.MarketData.CurrentPrice.Eur,
		coin.MarketData.CurrentPrice.Gbp,
		coin.MarketData.CurrentPrice.Inr,
	)
}
