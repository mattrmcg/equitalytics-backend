package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/buger/jsonparser"

	"github.com/mattrmcg/equitalytics-backend/config"
	"github.com/mattrmcg/equitalytics-backend/internal/db"
	"github.com/mattrmcg/equitalytics-backend/internal/services/data"
)

// updates the market price facts
func updateMarketPrice() {

	dbPool, err := db.CreateDBPool(config.Envs.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.CloseDBPool(dbPool)

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("unable to establish connection to database: %v\n", err)
	}
	log.Println("Connection to database established")

	dataService := data.NewDataService(dbPool)

	// retrieve slice of all companies from database
	companies, err := dataService.RetrieveCompanyMarketData(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	totalCounter := 0
	successCounter := 0
	startTime := time.Now()

	// iterate through companies
	for _, company := range companies {
		fmt.Printf("START CIK %v: %v\n", company.CIK, company.Ticker)

		totalCounter += 1

		// fetch market price
		marketPrice, err := fetchMarketPrice(company.Ticker)
		if err != nil {
			log.Printf("Error fetching market price: %v\n", err)
			fmt.Printf("END OF %v\n", company.Ticker)
			continue
		}

		// TODO: write calculations
		priceToEarningsRatio := calcPriceOverEarnings(marketPrice, company.EarningsPerShare)
		priceToBookRatio := calcPriceToBookRatio(marketPrice, company.BookValuePerShare)
		priceToSalesRatio := calcPriceToSalesRatio(marketPrice, company.RevenuePerShare)

		err = dataService.UpdateMarketPriceFacts(context.Background(), company.CIK, marketPrice, priceToEarningsRatio, priceToBookRatio, priceToSalesRatio)
		if err != nil {
			log.Printf("Error updating market price; %v\n", err)
		}

		successCounter += 1

		log.Println("Success!")

		fmt.Printf("END OF CIK%v: %v\n", company.CIK, company.Ticker)

	}

	log.Printf("\nTotal Companies: %v, Successful Updates: %v, Time elapsed: %v\n", totalCounter, successCounter, time.Since(startTime))

}

func fetchMarketPrice(ticker string) (float64, error) {
	url := fmt.Sprint("https://query1.finance.yahoo.com/v8/finance/chart/", ticker)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0.0, err
	}

	time.Sleep(time.Millisecond * 200)
	resp, err := client.Do(req)
	if err != nil {
		return 0.0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, err
	}

	marketPrice, err := jsonparser.GetFloat(body, "chart", "result", "[0]", "meta", "previousClose")
	if err != nil {
		return 0.0, err
	}
	return marketPrice, nil
}

func calcPriceOverEarnings(marketPrice float64, earningsPerShare float64) float64 {
	if earningsPerShare != 0 {
		return marketPrice / earningsPerShare
	}
	return float64(0.0)
}

func calcPriceToBookRatio(marketPrice float64, bookValuePerShare float64) float64 {
	if bookValuePerShare != 0 {
		return marketPrice / bookValuePerShare
	}
	return float64(0.0)
}

func calcPriceToSalesRatio(marketPrice float64, revenuePerShare float64) float64 {
	if revenuePerShare != 0 {
		return marketPrice / revenuePerShare
	}
	return float64(0.0)
}
