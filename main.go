package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
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

const PriceUrl = "https://api.coingecko.com/api/v3/coins/"
const DefaultCrypt = "bitcoin"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Which crypto price do you want to check? eg: etheruem,bitcoin,tether,usd-coin...etc")
	fmt.Printf("> ")
	scanner.Scan()
	crypto := strings.ToLower(scanner.Text())
	if crypto == "" {
		fmt.Println("You have not given an input. Check the current value of Bitcoin.")
		fmt.Println()
		priceCheck(DefaultCrypt)
	} else {
		fmt.Println("Please wait....")
		fmt.Println()
		priceCheck(crypto)
	}
}

func priceCheck(c string) {
	crypto := strings.ReplaceAll(strings.Trim(c, " "), " ", "-")
	currentTime := time.Now()
	res, err := http.Get(PriceUrl + crypto)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == 200 {
		var coin CoinGeckoData

		err2 := json.Unmarshal(body, &coin)

		if err2 != nil {
			log.Fatal(err2)
		}

		fmt.Printf("%s (%s) Price Today (%v %v, %v | %v): \n", coin.Name, strings.ToUpper(coin.Symbol), currentTime.Month(), currentTime.Day(), currentTime.Year(), currentTime.Format(time.Kitchen))
		fmt.Printf("USD: %.2f\nEUR: %.2f\nGBP: %.2f\nINR: %.2f\n", math.Round(coin.MarketData.CurrentPrice.Usd*100)/100, math.Round(coin.MarketData.CurrentPrice.Eur*100)/100, math.Round(coin.MarketData.CurrentPrice.Gbp*100)/100, math.Round(coin.MarketData.CurrentPrice.Inr*100)/100)
	} else {
		fmt.Printf("Crypto not found! Please check if \"%s\" is an valid crypto.\n", cases.Title(language.English).String(c))
	}
}
