package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/mattrmcg/equitalytics-backend/internal/models"
)

// Two approaches:
// 1) Use JSON Facts API with CIK list JSON

// OPTION 1 WORKFLOW

// Loop through CIK JSON
// For Each CIK
//		(if CIK not yet in database) (if it is, continue to next iteration)
//		Ping Submissions JSON for SIC, Exchange, Tickers, etc
//		Ping Facts JSON for company facts
// 		Unmarshal needed facts into struct
//		Build CompanyInfo struct with previous struct fields
//		Add company_info row into database

// CIKList struct is for unmarshalling company_tickers_exchange.json

var client *http.Client = &http.Client{}

// command line arg to trigger seed function
func main() {
	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "seed" {
		seed()
	} else {
		log.Fatal("unrecognized argument")
	}
}

func seed() {

	// var companies []models.CompanyInfo // for holding CompanyInfo objects before storing them in database
	//companiesList := list.New()
	var errors []error // for storing non-fatal errors, database will only be accessed if there are no stored errors

	// grab cik list json
	rawData, err := fetchCIKList()
	if err != nil {
		log.Fatal(err)
	}

	// marshal cik list json into CIKList struct
	cikList, err := unmarshalIntoCIKList(rawData)
	if err != nil {
		log.Fatal(err)
	}

	// iterate through CIKList data field to grab each CIK
	for _, n := range cikList.Data {
		// if cik not already stored in database
		if len(n) > 0 {
			cikInt, err := convertCIKToInt64(n[0])
			if err != nil {
				log.Fatal(err)
			}

			// NOTE: The reason we fetch and unmarshal the submissions json is because both the CIK list JSON and
			// the Facts data JSON don't contain important company information like SIC code, Exchanges, etc.
			sub, err := getSubmissionsWithCIK(cikInt)
			if err != nil {
				errors = append(errors, err)
				log.Println(err)
			}

			// any CIK that isn't a large accelerated filer will be returned as nil
			if sub != nil {
				// // Temp
				// printStr := fmt.Sprint(sub.Name, " : ", sub.Tickers[0], " : ", sub.SICDescription)
				// fmt.Println(printStr)
				// // Temp

				// Get Company Facts
				/*facts*/
				facts, err := getFactsWithCIK(cikInt)
				if err != nil {
					errors = append(errors, err)
					log.Println(err)
				}

				// verifyFactsStruct returns a list of errors for each fact that isn't populated with data
				verifyErrorsList := verifyFactsStruct(facts)
				// append the returned list of errors to our previously declared errors list using a spread operator
				errors = append(errors, verifyErrorsList...)

				info, err := fillCompanyInfoStruct(cikInt, sub, facts)
				// CHANGE
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(info, errors) // REMOVE THIS LINE

			}

		}
		//time.Sleep(250 * time.Millisecond)
	}

	log.Println("Success! Database has been seeded via the SEC API's!")
}

// fetch and return CIK List json
// returns byte slice contaning json data
func fetchCIKList() ([]byte, error) {

	// using company_tickers_exchange because format of regular company_tickers is bad
	req, err := http.NewRequest("GET", "https://www.sec.gov/files/company_tickers_exchange.json", nil)
	if err != nil {
		return nil, err
	}

	// User-Agent header is needed for SEC API's to accept requests
	req.Header.Add("User-Agent", "Matt McGuire matthewrmcguire56@gmail.com")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// unmarshals CIK request data into CIKList struct
func unmarshalIntoCIKList(rawData []byte) (*CIKList, error) {
	var cikList CIKList
	err := json.Unmarshal(rawData, &cikList)
	if err != nil {
		return nil, err
	}
	return &cikList, nil
}

// unmarshals submissions request data into Submissions struct
func unmarshalIntoSubmissions(rawData []byte) (*Submissions, error) {
	if !checkEntityTypeAndCategory(rawData) {
		return nil, nil
	}
	var sub Submissions
	err := json.Unmarshal(rawData, &sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func unmarshalIntoFacts(rawData []byte) (*Facts, error) {
	var facts Facts
	err := json.Unmarshal(rawData, &facts)
	if err != nil {
		return nil, err
	}

	return &facts, nil
}

// verifies that the cik is of entityType:"operating" and category"Large accelerated filer"
func checkEntityTypeAndCategory(rawData []byte) bool {
	entityType, err := jsonparser.GetString(rawData, "entityType")
	if err != nil {
		log.Fatalf("Couldn't parse entityType: %v", err)
	}

	category, err := jsonparser.GetString(rawData, "category")
	if err != nil {
		log.Fatalf("Couldn't parse category: %v", err)
	}

	if entityType != "operating" || category != "Large accelerated filer" {
		return false
	}
	return true
}

// converts cik float to type int64
func convertCIKToInt64(rawCIK any) (int64, error) {
	if cikFloat, ok := rawCIK.(float64); ok {
		cikStr := strconv.FormatFloat(cikFloat, 'f', -1, 64)
		cikInt, err := strconv.ParseInt(cikStr, 10, 64)
		if err != nil {
			log.Fatalf("Error converting CIK string to int64: %v", err)
		}

		return cikInt, nil
	} else {
		return 0, errors.New("cik was not of type float64")
	}
}

// fetchs submissions json using cik argument
func getSubmissionsWithCIK(cik int64) (*Submissions, error) {
	// the SEC API has a rate limit of 10 requests per second
	// the time.sleep ensures there is a wait period before each request to submissions
	// this makes it impossible to go over rate limit
	time.Sleep(150 * time.Millisecond)
	url := buildSubmissionsURL(cik)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %v", err)
	}

	req.Header.Add("User-Agent", "Matt McGuire matthewrmcguire56@gmail.com")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't make request: %v", err)
	}

	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read submissions body: %v", err)
	}

	sub, err := unmarshalIntoSubmissions(rawData)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal CIK %v into submissions struct: %v", cik, err)
	}

	return sub, nil
}

// builds the url that is used to fetch submissions json
func buildSubmissionsURL(cik int64) string {
	cikStr := convertCIKToURLString(cik)

	url := fmt.Sprint("https://data.sec.gov/submissions/CIK", cikStr, ".json")
	return url
}

func getFactsWithCIK(cik int64) (*Facts, error) {
	time.Sleep(150 * time.Millisecond)

	url := buildFactsURL(cik)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to read submissions body: %v", err)
	}

	req.Header.Add("User-Agent", "Matt McGuire matthewrmcguire56@gmail.com")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't make request: %v", err)
	}

	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read facts body: %v", err)
	}

	facts, err := unmarshalIntoFacts(rawData)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal cik %v into facts struct: %v", cik, err)
	}

	return facts, nil
}

func buildFactsURL(cik int64) string {
	cikStr := convertCIKToURLString(cik)

	url := fmt.Sprint("https://data.sec.gov/api/xbrl/companyfacts/CIK", cikStr, ".json")
	return url
}

// converts a cik integer to a string suitable for use in retrieval URLs
func convertCIKToURLString(cik int64) string {
	cikStr := strconv.FormatInt(cik, 10)
	for len(cikStr) < 10 {
		cikStr = fmt.Sprint("0", cikStr)
	}

	return cikStr
}

// This function verifies that each Fact retrieved actually contains data
// For each Fact that doesn't have its respective data, an error will be added to an error slice and logged.
// The error slice is returned so that it can be appended to the error slice in the seed function.
func verifyFactsStruct(facts *Facts) []error {
	var errList []error

	// We check that the length of the unit type struct array is not equal to 0 in order to verify that it's populated with data
	length := len(facts.Facts.USGAAP.NetIncomeLoss.Units.USD)
	if length == 0 {
		err := fmt.Errorf("NetIncomeLoss data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.Assets.Units.USD)
	if length == 0 {
		err := fmt.Errorf("assets data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.NetCashProvidedByUsedInOperatingActivities.Units.USD)
	if length == 0 {
		err := fmt.Errorf("NetCashProvidedByUsedInOperatingActivities data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.LongTermDebt.Units.USD)
	if length == 0 {
		err := fmt.Errorf("LongTermDebt data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.AssetsCurrent.Units.USD)
	if length == 0 {
		err := fmt.Errorf("AssetsCurrent data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.LiabilitiesCurrent.Units.USD)
	if length == 0 {
		err := fmt.Errorf("LiabiltiesCurrent data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.CommonStockSharesIssued.Units.Shares)
	if length == 0 {
		err := fmt.Errorf("CommonStockSharesIssued data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.Revenues.Units.USD)
	if length == 0 {
		err := fmt.Errorf("'Revenues' data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.CostOfGoodsSold.Units.USD)
	if length == 0 {
		err := fmt.Errorf("CostOfGoodsSold data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.OperatingIncomeLoss.Units.USD)
	if length == 0 {
		err := fmt.Errorf("OperatingIncomeLoss data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.StockholdersEquity.Units.USD)
	if length == 0 {
		err := fmt.Errorf("StockholdersEquity data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.CashAndCashEquivalentsAtCarryingValue.Units.USD)
	if length == 0 {
		err := fmt.Errorf("CashAndCashEquivalentsAtCarryingValue data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.ShortTermInvestments.Units.USD)
	if length == 0 {
		err := fmt.Errorf("ShortTermInvestments data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.AccountsReceivableNetCurrent.Units.USD)
	if length == 0 {
		err := fmt.Errorf("AccountsReceivableNetCurrent data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.InterestExpense.Units.USD)
	if length == 0 {
		err := fmt.Errorf("InterestExpense data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares)
	if length == 0 {
		err := fmt.Errorf("EarningsPerShareBasic data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.WeightedAverageNumberOfSharesOutstandingBasic.Units.Shares)
	if length == 0 {
		err := fmt.Errorf("WeightedAverageNumberOfSharesOutstandingBasic data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	length = len(facts.Facts.USGAAP.GrossProfit.Units.USD)
	if length == 0 {
		err := fmt.Errorf("GrossProfit data not present: %v", facts.EntityName)
		errList = append(errList, err)
		log.Println(err)
	}

	return errList

}

// creates and populates a CompanyInfo struct with data from Facts and Submissions structs
func fillCompanyInfoStruct(cik int64, sub *Submissions, facts *Facts) (*models.CompanyInfo, error) {
	// slice for collecting errors to be returned at the end
	var errors []error

	sic, err := strconv.ParseInt(sub.SIC, 10, 64)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}
	name, ticker, exchanges := sub.Name, sub.Tickers[0], sub.Exchanges

	// We need the three most recent yearly assets values later on for the Asset Turnover Ratio
	// However, we are only storing the latestAssets value in our CompanyInfo struct
	latestAssets, secondLatestAssets, thirdLatestAssets, err := getAssets(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	liabilities, err := getLiabilities(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	revenues, previousRevenues, err := getCurrentAndPreviousRevenues(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	costOfGoodsSold, err := getCostOfGoodsSold(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	grossProfit, previousGrossProfit, err := getCurrentAndPreviousGrossProfit(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	operatingIncomeLoss, err := getOperatingIncomeLoss(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	stockholdersEquity, previousStockholdersEquity, err := getCurrentAndPreviousYearStockholdersEquity(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	cashAndCashEquivalents, err := getCashAndCashEquivalents(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	shortTermInvestments, err := getShortTermInvestments(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	accountsReceivableNetCurrent, previousAccountsReceivableNetCurrent, err := getCurrentAndPreviousYearAccountsReceivableNetCurrent(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	interestExpense, err := getInterestExpense(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	weightedAverageNumberOfSharesOutstanding, err := getWeightedAverageNumberOfSharesOutstanding(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	bookValuePerShare := calcBookValuePerShare(stockholdersEquity, weightedAverageNumberOfSharesOutstanding)

	revenuePerShare := calcRevenuePerShare(revenues, weightedAverageNumberOfSharesOutstanding)

	// We first need to check if commonStockDividendsPerShareDeclared is present in the Facts struct our data was unmarshalled into
	var commonStockDividendsPerShareDeclared float64 = 0.0

	if checkIfCommonStockDividendsPerShareDeclaredExists(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared) {
		// If it exists, we grab it from the Facts struct using its associated retrieval function
		commonStockDividendsPerShareDeclared, err = getCommonStockDividendsPerShareDeclared(facts)
		if err != nil {
			log.Println(err)
			errors = append(errors, err)
		}
	}
	// If it the 'if' statement doesn't execute, then commonStockDividendsPerShareDeclared is initialized in our CompanyInfo
	// struct with a value of 0.0

	piotroskiScore := 0
	pointsInProfitability := 0
	pointsInLeverageLiquiditySourceOfFunds := 0
	pointsInOperatingEfficiency := 0

	// PIOTROSKI PROFITABILITY
	netIncome, err := getNetIncomeLoss(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}
	isNetIncomePositive := netIncome > 0

	returnOnAssets := calcReturnOnAssets(netIncome, latestAssets)
	isReturnOnAssetsPositive := returnOnAssets > 0

	operatingCashFlow, err := getOperatingCashFlow(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}
	isOperatingCashFlowPositive := operatingCashFlow > 0

	isOCFGreaterThanNetIncome := operatingCashFlow > netIncome

	if isNetIncomePositive {
		pointsInProfitability += 1
	}
	if isReturnOnAssetsPositive {
		pointsInProfitability += 1
	}
	if isOperatingCashFlowPositive {
		pointsInProfitability += 1
	}
	if isOCFGreaterThanNetIncome {
		pointsInProfitability += 1
	}

	piotroskiScore += pointsInProfitability

	// PIOTROSKI LEVERAGE LIQUIDITY SOURCE OF FUNDS

	// Point 5
	longTermDebt, previousLongTermDebt, err := getCurrentAndPreviousYearLongTermDebt(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}
	isCurrentLTDGreaterThanPreviousLTD := longTermDebt > previousLongTermDebt

	// Point 6
	assetsCurrent, previousAssetsCurrent, err := getCurrentAndPreviousYearAssetsCurrent(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	liabilitiesCurrent, previousLiabilitiesCurrent, err := getCurrentAndPreviousYearLiabilitiesCurrent(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	currentRatio := calcCurrentRatio(assetsCurrent, liabilitiesCurrent)
	previousCurrentRatio := calcCurrentRatio(previousAssetsCurrent, previousLiabilitiesCurrent)

	isCurrentCRGreaterThanPreviousCR := currentRatio > previousCurrentRatio

	// Point 7
	commonStockSharesIssued, previousCommonStockSharesIssued, err := getCurrentAndPreviousYearCommonStockSharesIssued(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}
	sharesIssuedInTheLastYear := commonStockSharesIssued - previousCommonStockSharesIssued
	noNewSharesIssued := sharesIssuedInTheLastYear <= 0

	if isCurrentLTDGreaterThanPreviousLTD {
		pointsInLeverageLiquiditySourceOfFunds += 1
	}
	if isCurrentCRGreaterThanPreviousCR {
		pointsInLeverageLiquiditySourceOfFunds += 1
	}
	if noNewSharesIssued {
		pointsInLeverageLiquiditySourceOfFunds += 1
	}

	piotroskiScore += pointsInLeverageLiquiditySourceOfFunds

	// PIOTROSKI OPERATING EFFICIENCY
	// Point 8
	grossProfitMargin := calcGrossProfitMargin(grossProfit, revenues)
	previousGrossProfitMargin := calcGrossProfitMargin(previousGrossProfit, previousRevenues)
	currentGPMGreaterThanPreviousGPM := grossProfitMargin > previousGrossProfitMargin

	// Point 9
	assetTurnoverRatio := calcAssetTurnoverRatio(revenues, latestAssets, secondLatestAssets)
	previousAssetTurnoverRatio := calcAssetTurnoverRatio(previousRevenues, secondLatestAssets, thirdLatestAssets)
	isCurrentATRGreaterThanPreviousATR := assetTurnoverRatio > previousAssetTurnoverRatio

	if currentGPMGreaterThanPreviousGPM {
		pointsInOperatingEfficiency += 1
	}
	if isCurrentATRGreaterThanPreviousATR {
		pointsInOperatingEfficiency += 1
	}

	piotroskiScore += pointsInOperatingEfficiency

	// DONE WITH PIOTROSKI
	// ADDITIONAL METRICS

	// PROFITABILITY
	operatingProfitMargin := calcOperatingProfitMargin(operatingIncomeLoss, revenues)
	netProfitMargin := calcNetProfitMargin(netIncome, revenues)
	returnOnEquity := calcReturnOnEquity(netIncome, stockholdersEquity, previousStockholdersEquity)

	// LIQUIDITY
	quickRatio := calcQuickRatio(cashAndCashEquivalents, shortTermInvestments, accountsReceivableNetCurrent, liabilitiesCurrent)

	// SOLVENCY
	debtToEquityRatio := calcDebtToEquityRatio(liabilities, stockholdersEquity)
	debtToAssetsRatio := calcDebtToAssetsRatio(liabilities, latestAssets)
	interestCoverageRatio := calcInterestCoverageRatio(operatingIncomeLoss, interestExpense)

	// EFFICIENCY
	// inventoryTurnoverRatio :=
	receivablesTurnoverRatio := calcReceivablesTurnoverRatio(revenues, accountsReceivableNetCurrent, previousAccountsReceivableNetCurrent)

	// VALUATION
	priceToEarningsRatio := 0.0
	priceToBookRatio := 0.0
	priceToSalesRatio := 0.0
	dividendYield := 0.0
	earningsPerShare, err := getEarningsPerShareBasic(facts)
	if err != nil {
		log.Println(err)
		errors = append(errors, err)
	}

	companyInfo := &models.CompanyInfo{
		CIK:         cik,
		SIC:         sic,
		CompanyName: name,
		Ticker:      ticker,
		Exchanges:   exchanges,

		Assets:                                   latestAssets,
		Liabilities:                              liabilities,
		Revenues:                                 revenues,
		PreviousYearRevenues:                     previousRevenues,
		CostOfGoodsSold:                          costOfGoodsSold,
		GrossProfit:                              grossProfit,
		PreviousYearGrossProfit:                  previousGrossProfit,
		OperatingIncomeLoss:                      operatingIncomeLoss,
		StockholdersEquity:                       stockholdersEquity,
		PreviousYearStockholdersEquity:           previousStockholdersEquity,
		CashAndCashEquivalents:                   cashAndCashEquivalents,
		ShortTermInvestments:                     shortTermInvestments,
		AccountsReceivableNetCurrent:             accountsReceivableNetCurrent,
		PreviousYearAccountsReceivableNetCurrent: previousAccountsReceivableNetCurrent,
		InterestExpense:                          interestExpense,
		WeightedAverageNumberOfSharesOutstanding: weightedAverageNumberOfSharesOutstanding,
		BookValuePerShare:                        bookValuePerShare,
		RevenuePerShare:                          revenuePerShare,
		CommonStockDividendsPerShareDeclared:     commonStockDividendsPerShareDeclared,

		// Piotroski
		PiotroskiScore:                         piotroskiScore,
		PointsInProfitability:                  pointsInProfitability,
		PointsInLeverageLiquiditySourceOfFunds: pointsInLeverageLiquiditySourceOfFunds,
		PointsInOperatingEfficiency:            pointsInOperatingEfficiency,

		// PIOTROSKI PROFITABILITIY
		// Point 1
		NetIncome:           netIncome,
		IsNetIncomePositive: isNetIncomePositive,

		// Point 2
		ReturnOnAssets:           returnOnAssets,
		IsReturnOnAssetsPositive: isReturnOnAssetsPositive,

		// Point 3
		OperatingCashFlow:           operatingCashFlow,
		IsOperatingCashFlowPositive: isOperatingCashFlowPositive,

		// Point 4
		IsOCFGreaterThanNetIncome: isOCFGreaterThanNetIncome,

		// PIOTROSKI LEVERAGE LIQUIDITY SOURCE OF FUNDS
		// Point 5
		LongTermDebt:                       longTermDebt,
		PreviousYearLongTermDebt:           previousLongTermDebt,
		IsCurrentLTDGreaterThanPreviousLTD: isCurrentLTDGreaterThanPreviousLTD,

		// Point 6
		AssetsCurrent:                    assetsCurrent,
		LiabilitiesCurrent:               liabilitiesCurrent,
		PreviousYearAssetsCurrent:        previousAssetsCurrent,
		PreviousYearLiabilitiesCurrent:   previousLiabilitiesCurrent,
		CurrentRatio:                     currentRatio,
		PreviousYearCurrentRatio:         previousCurrentRatio,
		IsCurrentCRGreaterThanPreviousCR: isCurrentCRGreaterThanPreviousCR,

		// Point 7
		CommonStockSharesIssued:             commonStockSharesIssued,
		PreviousYearCommonStockSharesIssued: previousCommonStockSharesIssued,
		SharesIssuedInTheLastYear:           sharesIssuedInTheLastYear,
		NoNewSharesIssued:                   noNewSharesIssued,

		// PIOTROSKI OPERATING EFFICIENCY
		// Point 8
		GrossProfitMargin:                grossProfitMargin,
		PreviousYearGrossProfitMargin:    previousGrossProfitMargin,
		CurrentGPMGreaterThanPreviousGPM: currentGPMGreaterThanPreviousGPM,

		// Point 9
		AssetTurnoverRatio:                 assetTurnoverRatio,
		PreviousYearAssetTurnoverRatio:     previousAssetTurnoverRatio,
		IsCurrentATRGreaterThanPreviousATR: isCurrentATRGreaterThanPreviousATR,

		// DONE WITH PIOTROSKI SCORING, MOVE ON TO OTHER METRICS
		// PROFITABILITY
		OperatingProfitMargin: operatingProfitMargin,
		NetProfitMargin:       netProfitMargin,
		ReturnOnEquity:        returnOnEquity,

		// LIQUIDITY
		QuickRatio: quickRatio,

		// SOLVENCY
		DebtToEquityRatio:     debtToEquityRatio,
		DebtToAssetsRatio:     debtToAssetsRatio,
		InterestCoverageRatio: interestCoverageRatio,

		// EFFICIENCY
		ReceivablesTurnoverRatio: receivablesTurnoverRatio,

		// VALUATION
		PriceToEarningsRatio: priceToEarningsRatio,
		PriceToBookRatio:     priceToBookRatio,
		PriceToSalesRatio:    priceToSalesRatio,
		DividendYield:        dividendYield,
		EarningsPerShare:     earningsPerShare,

		// TIMESTAMPS
		FilingDataUpdatedAt: time.Now(),
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("CompanyInfo struct couldn't be filled for CIK %v", cik)
	}

	return companyInfo, nil
}

// ##############################################################################################################################
// # FUNCTIONS FOR RETRIEVING YEARLY FILING VALUES FROM FACTS STRUCT															#
// ##############################################################################################################################

// get latest yearly NetIncomeLoss value from struct
func getNetIncomeLoss(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.NetIncomeLoss.Units.USD) - 1
	netIncomeLossValueIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.NetIncomeLoss, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getNetIncomeLosss() returned an error: %v", err)
	}

	netIncomeLossValue := facts.Facts.USGAAP.NetIncomeLoss.Units.USD[netIncomeLossValueIndex].Val

	intVal, err := convertJSONNumberToInt64(netIncomeLossValue)
	if err != nil {
		return 0, fmt.Errorf("getNetIncomeLoss(): couldn't convert to int64 value: %v", err)
	}

	return intVal, nil
}

// getAssets() retrieves the three latest yearly reported Assets values from the facts struct that the JSON data was unmarshalled into.
// We need the three latest values in order to calculate the Current and Previous year Asset Turnover Ratio later on.
func getAssets(facts *Facts) (int64, int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.Assets.Units.USD) - 1
	latestAssetsIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.Assets, startingIndex)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("getAssets() returned an error: %v", err)
	}

	secondLatestAssetsIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.Assets, latestAssetsIndex-1)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("getAssets() returned an error: %v", err)
	}

	thirdLatestAssetsIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.Assets, secondLatestAssetsIndex-1)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("getAssets() returned an error: %v", err)
	}

	latestAssetsValue := facts.Facts.USGAAP.Assets.Units.USD[latestAssetsIndex].Val
	secondLatestAssetsValue := facts.Facts.USGAAP.Assets.Units.USD[secondLatestAssetsIndex].Val
	thirdLatestAssetsValue := facts.Facts.USGAAP.Assets.Units.USD[thirdLatestAssetsIndex].Val

	intVal, err := convertJSONNumberToInt64(latestAssetsValue)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("getAssets(): couldn't convert to int64 value: %v", err)
	}

	secondIntVal, err := convertJSONNumberToInt64(secondLatestAssetsValue)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("getAssets(): couldn't convert to int64 value: %v", err)
	}

	thirdIntVal, err := convertJSONNumberToInt64(thirdLatestAssetsValue)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("getAssets(): couldn't convert to int64 value: %v", err)
	}

	return intVal, secondIntVal, thirdIntVal, nil
}

// get latest yearly Liabilities value from struct
func getLiabilities(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.Liabilities.Units.USD) - 1
	latestLiabilitiesValueIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.Liabilities, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getLiabilities() returned an error: %v", err)
	}

	latestLiabilitiesValue := facts.Facts.USGAAP.Liabilities.Units.USD[latestLiabilitiesValueIndex].Val

	intVal, err := convertJSONNumberToInt64(latestLiabilitiesValue)
	if err != nil {
		return 0, fmt.Errorf("getLiabilities(): couldn't convert to int64: %v", err)
	}

	return intVal, nil
}

// get latest yearly OperatingCashFlow value
func getOperatingCashFlow(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.NetCashProvidedByUsedInOperatingActivities.Units.USD) - 1
	operatingCashFlowValueIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.NetCashProvidedByUsedInOperatingActivities, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getOperatingCashFlow() returned an error: %v", err)
	}

	operatingCashFlowValue := facts.Facts.USGAAP.NetCashProvidedByUsedInOperatingActivities.Units.USD[operatingCashFlowValueIndex].Val

	intVal, err := convertJSONNumberToInt64(operatingCashFlowValue)
	if err != nil {
		return 0, fmt.Errorf("getOperatingCashFlow(): couldn't convert to int64 value: %v", err)
	}

	return intVal, nil
}

// get latest and previous yearly LongTermDebt
func getCurrentAndPreviousYearLongTermDebt(facts *Facts) (int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.LongTermDebt.Units.USD) - 1
	latestLongTermDebtValueIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.LongTermDebt, startingIndex)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearLongTermDebt() returned an error: %v", err)
	}

	previousLongTermDebtValueIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.LongTermDebt, latestLongTermDebtValueIndex-1)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearLongTermDebt() returned an error: %v", err)
	}
	latestLongTermDebtValue := facts.Facts.USGAAP.LongTermDebt.Units.USD[latestLongTermDebtValueIndex].Val
	previousLongTermDebtValue := facts.Facts.USGAAP.LongTermDebt.Units.USD[previousLongTermDebtValueIndex].Val

	latestIntVal, err := convertJSONNumberToInt64(latestLongTermDebtValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearLongTermDebt(): couldn't convert to int64 value: %v", err)
	}

	previousIntVal, err := convertJSONNumberToInt64(previousLongTermDebtValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearLongTermDebt(): couldn't convert to int64 value: %v", err)
	}

	return latestIntVal, previousIntVal, nil
}

// get latest and previous yearly AssetsCurrent
func getCurrentAndPreviousYearAssetsCurrent(facts *Facts) (int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.AssetsCurrent.Units.USD) - 1
	latestAssetsCurrentIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.AssetsCurrent, startingIndex)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearAssetsCurrent() returned an error: %v", err)
	}

	previousAssetsCurrentIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.AssetsCurrent, latestAssetsCurrentIndex-1)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearAssetsCurrent() returned an error: %v", err)
	}
	latestAssetsCurrentValue := facts.Facts.USGAAP.AssetsCurrent.Units.USD[latestAssetsCurrentIndex].Val
	previousAssetsCurrentValue := facts.Facts.USGAAP.AssetsCurrent.Units.USD[previousAssetsCurrentIndex].Val

	latestIntVal, err := convertJSONNumberToInt64(latestAssetsCurrentValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousAssetsCurrent(): couldn't convert to int64 value: %v", err)
	}

	previousIntVal, err := convertJSONNumberToInt64(previousAssetsCurrentValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearAssetsCurrent(): couldn't convert to int64 value: %v", err)
	}

	return latestIntVal, previousIntVal, nil
}

// get latest and previous yearly LiabilitiesCurrent
func getCurrentAndPreviousYearLiabilitiesCurrent(facts *Facts) (int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.LiabilitiesCurrent.Units.USD) - 1
	latestLiabilitiesCurrentIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.LiabilitiesCurrent, startingIndex)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearLiabilitiesCurrent() returned an error: %v", err)
	}

	previousLiabilitiesCurrentIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.LiabilitiesCurrent, latestLiabilitiesCurrentIndex-1)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearLiabilitiesCurrent() returned an error: %v", err)
	}
	latestLiabilitiesCurrentValue := facts.Facts.USGAAP.LiabilitiesCurrent.Units.USD[latestLiabilitiesCurrentIndex].Val
	previousLiabilitiesCurrentValue := facts.Facts.USGAAP.LiabilitiesCurrent.Units.USD[previousLiabilitiesCurrentIndex].Val

	latestIntVal, err := convertJSONNumberToInt64(latestLiabilitiesCurrentValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearLiabilitiesCurrent() returned an error: %v", err)
	}

	previousIntVal, err := convertJSONNumberToInt64(previousLiabilitiesCurrentValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearLiabilitiesCurrent() returned an error: %v", err)
	}

	return latestIntVal, previousIntVal, nil
}

func getCurrentAndPreviousYearCommonStockSharesIssued(facts *Facts) (int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.CommonStockSharesIssued.Units.Shares) - 1
	latestCommonStockSharesIssuedIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.CommonStockSharesIssued, startingIndex)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearCommonStockSharesIssued() returned an error: %v", err)
	}

	previousCommonStockSharesIssuedIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.CommonStockSharesIssued, latestCommonStockSharesIssuedIndex-1)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearCommonStockSharesIssued() returned an error: %v", err)
	}
	latestCommonStockSharesIssuedValue := facts.Facts.USGAAP.CommonStockSharesIssued.Units.Shares[latestCommonStockSharesIssuedIndex].Val
	previousCommonStockSharesIssuedValue := facts.Facts.USGAAP.CommonStockSharesIssued.Units.Shares[previousCommonStockSharesIssuedIndex].Val

	latestIntVal, err := convertJSONNumberToInt64(latestCommonStockSharesIssuedValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearCommonStockSharesIssued() couldn't convert to int64: %v", err)
	}

	previousIntVal, err := convertJSONNumberToInt64(previousCommonStockSharesIssuedValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearCommonStockSharesIssued() couln't convert to int64: %v", err)
	}

	return latestIntVal, previousIntVal, nil
}

// get latest yearly Revenues value from struct
func getCurrentAndPreviousRevenues(facts *Facts) (int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.Revenues.Units.USD) - 1
	revenuesIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.Revenues, startingIndex)
	if err != nil {
		return 0, 0, fmt.Errorf("getRevenues() returned an error: %v", err)
	}

	previousRevenuesIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.Revenues, revenuesIndex-1)
	if err != nil {
		return 0, 0, fmt.Errorf("getRevenues() returned an error: %v", err)
	}

	revenuesValue := facts.Facts.USGAAP.Revenues.Units.USD[revenuesIndex].Val
	previousRevenuesValue := facts.Facts.USGAAP.Revenues.Units.USD[previousRevenuesIndex].Val

	intVal, err := convertJSONNumberToInt64(revenuesValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getRevenues(): couldn't convert to int64: %v", err)
	}

	previousIntVal, err := convertJSONNumberToInt64(previousRevenuesValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getRevenues(): couldn't convert to int64: %v", err)
	}

	return intVal, previousIntVal, nil
}

func getCostOfGoodsSold(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.CostOfGoodsSold.Units.USD) - 1
	cogsValueIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.CostOfGoodsSold, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getCostOfGoodsSold() returned an error: %v", err)
	}

	cogsValue := facts.Facts.USGAAP.Revenues.Units.USD[cogsValueIndex].Val

	intVal, err := convertJSONNumberToInt64(cogsValue)
	if err != nil {
		return 0, fmt.Errorf("getCostOfGoodsSold(): couldn't convert to int64: %v", err)
	}

	return intVal, nil
}

func getOperatingIncomeLoss(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.OperatingIncomeLoss.Units.USD) - 1
	operatingIncomeLossIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.OperatingIncomeLoss, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getOperatingIncomeLoss() returned an error: %v", err)
	}

	operatingIncomeLossValue := facts.Facts.USGAAP.OperatingIncomeLoss.Units.USD[operatingIncomeLossIndex].Val

	intVal, err := convertJSONNumberToInt64(operatingIncomeLossValue)
	if err != nil {
		return 0, fmt.Errorf("getOperatingIncomeLoss(): couldn't convert to int64: %v", err)
	}

	return intVal, nil
}

// StockholdersEquity + PreviousYearStockholdersEquity
func getCurrentAndPreviousYearStockholdersEquity(facts *Facts) (int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.StockholdersEquity.Units.USD) - 1
	latestStockholdersEquityIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.StockholdersEquity, startingIndex)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearStockholdersEquity() returned an error: %v", err)
	}

	previousStockholdersEquityIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.StockholdersEquity, latestStockholdersEquityIndex-1)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearStockholdersEquity() returned an error: %v", err)
	}
	latestStockholdersEquityValue := facts.Facts.USGAAP.StockholdersEquity.Units.USD[latestStockholdersEquityIndex].Val
	previousStockholdersEquityValue := facts.Facts.USGAAP.StockholdersEquity.Units.USD[previousStockholdersEquityIndex].Val

	latestIntVal, err := convertJSONNumberToInt64(latestStockholdersEquityValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearStockholdersEquity() couldn't convert to int64: %v", err)
	}

	previousIntVal, err := convertJSONNumberToInt64(previousStockholdersEquityValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearStockholdersEquity() couln't convert to int64: %v", err)
	}

	return latestIntVal, previousIntVal, nil
}

func getCashAndCashEquivalents(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.CashAndCashEquivalentsAtCarryingValue.Units.USD) - 1
	cashAndCashEquivalentsIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.CashAndCashEquivalentsAtCarryingValue, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getCashAndCashEquivalents() returned an error: %v", err)
	}

	cashAndCashEqivalentsValue := facts.Facts.USGAAP.CashAndCashEquivalentsAtCarryingValue.Units.USD[cashAndCashEquivalentsIndex].Val

	intVal, err := convertJSONNumberToInt64(cashAndCashEqivalentsValue)
	if err != nil {
		return 0, fmt.Errorf("getCashAndCashEquivalents(): couldn't convert to int64: %v", err)
	}

	return intVal, nil
}

func getShortTermInvestments(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.ShortTermInvestments.Units.USD) - 1
	shortTermInvestmentsIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.ShortTermInvestments, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getShortTermInvestments() returned an error: %v", err)
	}

	shortTermInvestmentsValue := facts.Facts.USGAAP.ShortTermInvestments.Units.USD[shortTermInvestmentsIndex].Val

	intVal, err := convertJSONNumberToInt64(shortTermInvestmentsValue)
	if err != nil {
		return 0, fmt.Errorf("getOperatingIncomeLoss(): couldn't convert to int64: %v", err)
	}

	return intVal, nil
}

func getCurrentAndPreviousYearAccountsReceivableNetCurrent(facts *Facts) (int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.AccountsReceivableNetCurrent.Units.USD) - 1
	latestARNCIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.AccountsReceivableNetCurrent, startingIndex)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearAccountsReceivableNetCurrent() returned an error: %v", err)
	}

	previousARNCIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.AccountsReceivableNetCurrent, latestARNCIndex-1)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearAccountsReceivableNetCurrent() returned an error: %v", err)
	}
	latestARNCValue := facts.Facts.USGAAP.AccountsReceivableNetCurrent.Units.USD[latestARNCIndex].Val
	previousARNCValue := facts.Facts.USGAAP.AccountsReceivableNetCurrent.Units.USD[previousARNCIndex].Val

	latestIntVal, err := convertJSONNumberToInt64(latestARNCValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearAccountsReceivableNetCurrent() couldn't convert to int64: %v", err)
	}

	previousIntVal, err := convertJSONNumberToInt64(previousARNCValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousYearAccountsReceivableNetCurrent() couldn't convert to int64: %v", err)
	}

	return latestIntVal, previousIntVal, nil
}

func getInterestExpense(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.InterestExpense.Units.USD) - 1
	interestExpenseIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.InterestExpense, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getInterestExpense() returned an error: %v", err)
	}

	interestExpenseValue := facts.Facts.USGAAP.InterestExpense.Units.USD[interestExpenseIndex].Val

	intVal, err := convertJSONNumberToInt64(interestExpenseValue)
	if err != nil {
		return 0, fmt.Errorf("getInterestExpense(): couldn't convert to int64: %v", err)
	}

	return intVal, nil
}

func getEarningsPerShareBasic(facts *Facts) (float64, error) {
	startingIndex := len(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares) - 1
	epsIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.EarningsPerShareBasic, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getEarningsPerShareBasic() returned an error: %v", err)
	}

	epsValue := facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares[epsIndex].Val

	floatVal, err := convertJSONNumberToFloat64(epsValue)
	if err != nil {
		return 0, fmt.Errorf("getEarningsPerShareBasic(): couldn't convert to float64: %v", err)
	}

	return floatVal, nil
}

// func getCurrentAndPreviousYearInventoryNet() {

// }

func getWeightedAverageNumberOfSharesOutstanding(facts *Facts) (int64, error) {
	startingIndex := len(facts.Facts.USGAAP.WeightedAverageNumberOfSharesOutstandingBasic.Units.Shares) - 1
	wansoIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.WeightedAverageNumberOfSharesOutstandingBasic, startingIndex)
	if err != nil {
		return 0, fmt.Errorf("getWeightedAverageNumberOfSharesOutsanding() returned an error: %v", err)
	}

	wansoValue := facts.Facts.USGAAP.WeightedAverageNumberOfSharesOutstandingBasic.Units.Shares[wansoIndex].Val

	intVal, err := convertJSONNumberToInt64(wansoValue)
	if err != nil {
		return 0, fmt.Errorf("getWeightedAverageNumberOfSharesOutstanding(): couldn't convert to int64: %v", err)
	}

	return intVal, nil
}

func getCurrentAndPreviousGrossProfit(facts *Facts) (int64, int64, error) {
	startingIndex := len(facts.Facts.USGAAP.GrossProfit.Units.USD) - 1
	grossProfitIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.GrossProfit, startingIndex)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousGrossProfit() returned an error: %v", err)
	}

	previousGrossProfitIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.GrossProfit, grossProfitIndex-1)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousGrossProfit() returned an error %v", err)
	}

	grossProfitValue := facts.Facts.USGAAP.GrossProfit.Units.USD[grossProfitIndex].Val
	previousGrossProfitValue := facts.Facts.USGAAP.GrossProfit.Units.USD[previousGrossProfitIndex].Val

	intVal, err := convertJSONNumberToInt64(grossProfitValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousGrossProfit(): couldn't convert to int64: %v", err)
	}

	previousIntVal, err := convertJSONNumberToInt64(previousGrossProfitValue)
	if err != nil {
		return 0, 0, fmt.Errorf("getCurrentAndPreviousGrossProfit(): couldn't convert to int64: %v", err)
	}

	return intVal, previousIntVal, nil
}

// MIGHT NOT EXIST - need to check if exists before retrieving
func getCommonStockDividendsPerShareDeclared(facts *Facts) (float64, error) {
	startingIndex := len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares)
	latestIndex, err := getNextYearlyFilingValueIndex(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared, startingIndex)
	if err != nil {
		return 0.0, fmt.Errorf("getCommonStockDividendsPerShareDeclared() returned an error: %v", err)
	}

	value := facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares[latestIndex].Val

	floatVal, err := convertJSONNumberToFloat64(value)
	if err != nil {
		return 0.0, fmt.Errorf("getCommonStockDividendsPerShareDeclared(): couldn't convert to float64: %v", err)
	}

	return floatVal, nil
}

// ##############################################################################################################################
// # FUNCTIONS FOR CACULATING ADDITIONAL FACTS AND RATIOS																		#
// ##############################################################################################################################

// calculate BookValuePerShare given StockholdersEquity and WeightedAverageNumberOfSharesOutstanding
func calcBookValuePerShare(stockholdersEquity int64, weightedAverageNumberOfSharesOutstanding int64) float64 {
	// StockholdersEquity / WeightedAverageNumberOfSharesOustanding
	return float64(stockholdersEquity) / float64(weightedAverageNumberOfSharesOutstanding)
}

// calculate RevenuePerShare given Revenues and WANSO
func calcRevenuePerShare(revenues int64, weightedAverageNumberOfSharesOutstanding int64) float64 {
	return float64(revenues) / float64(weightedAverageNumberOfSharesOutstanding)
}

// func calculatePiotroskiScore() {

// }

// Calculate the first 4 points of the piotroski score
// 1 - positive net income
// 2 - positive return on assets
// 3 - positive operating cash flow
// 4 - ocf greater than net income
// func calculatePiostroskiProfitability(netIncome int64, roa int64, ocf int64) int {
// total := 0
// 	// Calculate point 1
// 	if netIncome > 0 {
// 		total += 1
// 	}

// 	// Calculate point 2
// 	if roa > 0 {
// 		total += 1
// 	}

// 	// Calculate point 3
// 	if ocf > 0 {
// 		total += 1
// 	}

// 	if ocf > netIncome {
// 		total += 1
// 	}

// 	return total
// }

// func calculatePiotroskiLeverageLiquiditySourceOfFunds() {

// }

// func calculatePiotroskiOperatingEfficiency() {

// }

// DONT NEED
// func isNetIncomePositive(netIncome int64) bool {
// 	return netIncome >= 0
// }

func calcReturnOnAssets(netIncomeLoss int64, assets int64) float64 {
	return float64(netIncomeLoss) / float64(assets)
}

func calcCurrentRatio(assetsCurrent int64, liabilitiesCurrent int64) float64 {
	return float64(assetsCurrent) / float64(liabilitiesCurrent)
}

func calcGrossProfitMargin(grossProfit int64, revenues int64) float64 {
	return float64(grossProfit) / float64(revenues)
}

func calcAssetTurnoverRatio(revenues int64, currAssets int64, prevAssets int64) float64 {
	return float64(revenues) / ((float64(currAssets) + float64(prevAssets)) / 2.0)
}

func calcOperatingProfitMargin(operatingIncomeLoss int64, revenues int64) float64 {
	return float64(operatingIncomeLoss) / float64(revenues)
}

func calcNetProfitMargin(netIncomeLoss int64, revenues int64) float64 {
	return float64(netIncomeLoss) / float64(revenues)
}

func calcReturnOnEquity(netIncomeLoss int64, currStockholdersEquity int64, prevStockholdersEquity int64) float64 {
	diffStockholdersEquity := float64(currStockholdersEquity) - float64(prevStockholdersEquity)
	return float64(netIncomeLoss) / diffStockholdersEquity
}

func calcQuickRatio(cashAndCashEquivalents int64, shortTermInvestments int64, accountsReceivableNetCurrent int64, liabilitiesCurrent int64) float64 {
	return (float64(cashAndCashEquivalents) + float64(shortTermInvestments) + float64(accountsReceivableNetCurrent)) / float64(liabilitiesCurrent)
}

func calcDebtToEquityRatio(liabilities int64, stockholdersEquity int64) float64 {
	return float64(liabilities) / float64(stockholdersEquity)
}

func calcDebtToAssetsRatio(liabilities int64, assets int64) float64 {
	return float64(liabilities) / float64(assets)
}

func calcInterestCoverageRatio(operatingIncomeLoss int64, interestExpense int64) float64 {
	return float64(operatingIncomeLoss) / float64(interestExpense)
}

func calcReceivablesTurnoverRatio(revenues int64, currAccountsReceivableNetCurrent int64, prevAccountsReceivableNetCurrent int64) float64 {
	avgAccountsReceivable := float64(currAccountsReceivableNetCurrent+prevAccountsReceivableNetCurrent) / 2.0
	return float64(revenues) / avgAccountsReceivable
}

// NEXT THREE CALCULATIONS NEED MARKET PRICE
// func calcPriceOverEarnings() {

// }

// func calcPriceToBookRatio() {

// }

// func calcPriceToSalesRatio() {

// }

// type AnyFact interface {
// }

// getLatestYearlyFilingIndex() grabs the index of the unit object containing the latest yearly filing.
// This function is generic and can be used with any valid fact struct.
// func getLatestYearlyFilingIndex[T any](fact T) int {
// 	return 0
// }

// Generic function to get the latest yearly filing value for a given fact.
// This should work with any fact, so long as the fact struct itself is passed.
// i.e. FactsStruct.Facts.USGAAP.AccountsReceivable

// returns a boolean value indicating if the CommonStockDividendsPerShareDeclared fact exists
type CommonStockDividendsPerShareDeclared struct {
	Units struct {
		USDOverShares []FactUnit
	}
}

// checks the length of the CommonStockDividendsPerShareDeclared field
// if the length is 0, then the fact wasn't declared in the original JSON
func checkIfCommonStockDividendsPerShareDeclaredExists(fact struct {
	Units struct {
		USDOverShares []FactUnit "json:\"USD/shares\""
	} "json:\"units\""
}) bool {

	return len(fact.Units.USDOverShares) != 0
}

func getNextYearlyFilingValueIndex(fact interface{}, startIndex int) (int, error) {
	v := reflect.ValueOf(fact)

	// if fact object is not a struct
	if v.Kind() != reflect.Struct {

		return 0, fmt.Errorf("expected fact to be a struct")
	}

	// Access "Units" Field
	unitsField := v.FieldByName("Units")
	if !unitsField.IsValid() || unitsField.Kind() != reflect.Struct {
		return 0, fmt.Errorf("expected a units field of struct type")
	}

	var unitField reflect.Value
	// Determine what type of unit the value is stored as (USD, Shares, or USD/Shares)
	if (unitsField.FieldByName("USD")).IsValid() {
		unitField = unitsField.FieldByName("USD")
	} else if (unitsField.FieldByName("Shares")).IsValid() {
		unitField = unitsField.FieldByName("Shares")
	} else if (unitsField.FieldByName("USDOverShares")).IsValid() {
		unitField = unitsField.FieldByName("USDOverShares")
	} else {
		return 0, fmt.Errorf("expected a unitField of USD, Shares, or USDOverShares")
	}

	// the unitField should be a slice of structs
	if unitField.Kind() != reflect.Slice {
		return 0, fmt.Errorf("expected unitField to be of type slice")
	}

	for i := startIndex; i > 0; i-- {
		item := unitField.Index(i)

		if item.Kind() != reflect.Struct {
			log.Printf("Unit field at index %v is not a struct\n", i)
			continue
		}

		// initialize reflect values for each field needed to determine filing year
		fpField := item.FieldByName("FP")
		formField := item.FieldByName("Form")
		frameField := item.FieldByName("Frame")

		// declare bool variables to determine if all fields exist and are of string type
		allValid := fpField.IsValid() && formField.IsValid() && frameField.IsValid()
		allStrings := (fpField.Kind() == reflect.String) && (formField.Kind() == reflect.String) && (frameField.Kind() == reflect.String)

		// fmt.Printf("allValid: %v | allStrings: %v\n", allValid, allStrings)
		if allValid && allStrings {

			// Check if the current element represents a yearly filing
			// frame should be empty, filing period should be 'FY' (filing year), and form type should be '10-K'

			if (frameField.String() == "" || frameField.String() == "CY2023") && fpField.String() == "FY" && formField.String() == "10-K" {
				return i, nil
			}
		}

	}

	return 0, fmt.Errorf("no adequate filing year found")
}

// this function converts a given json.Number type to an int64
func convertJSONNumberToInt64(num json.Number) (int64, error) {
	// attempts to convert json.Number directly to type int64
	// if no error, return the converted int64 value
	if int64Val, err := num.Int64(); err == nil {
		return int64Val, nil
	}

	// if json.Number cannot be directly converted to int64, attempt to convert it to float64
	// if an error is returned, the json.Number is not suited for conversion to either int64 or float64
	float64Val, err := num.Float64()
	if err != nil {
		return 0, fmt.Errorf("failed to convert json.Number to float64: %v", err)
	}

	// round float64 value to an int64 value and return
	int64Val := int64(math.Round(float64Val))
	return int64Val, nil
}

func convertJSONNumberToFloat64(num json.Number) (float64, error) {
	float64Val, err := num.Float64()
	if err != nil {
		return 0.0, fmt.Errorf("failed to convert json.Number to float64: %v", err)
	}
	return float64Val, nil
}

// CURRENTLY NOT USING OPTION 2
// Haven't figured out how to extract XBRL data from inlineXBRL, probably need an html parser to do so

// 2) Use CIK list with Submissions JSON to download and parse xbrl filings
// Option 2 is much harder and more time consuming, but ensures a higher degree of data integrity and information

// OPTION 2 WORKFLOW

// Loop through CIK JSON
//
// For each CIK
// 		ping /submissions endpoint
//		get most recent 10 Q filing accession number and filename
//		use accession number and filename to construct URL to EDGAR archives
//		retrieve htm file
// 		use and html parser ??? to extract iXBRL data
// 		use XBRL parser to unmarshal needed facts into XBRL struct (cik, sic, ticker, EPS, etc.)
//		use XBRL struct and functions from data_service to create Info struct and store in db
//		rinse and repeat
