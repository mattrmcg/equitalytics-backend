package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

				info := fillCompanyInfoStruct(cikInt, sub, facts)

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
		err := fmt.Errorf("Assets data not present: %v", facts.EntityName)
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
		err := fmt.Errorf("Revenues data not present: %v", facts.EntityName)
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
		err := fmt.Errorf("NetCashProvidedByUsedInOperatingActivities data not present: %v", facts.EntityName)
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
func fillCompanyInfoStruct(cik int64, sub *Submissions, factsStruct *Facts) *models.CompanyInfo {
	sic, name, ticker, exchanges := sub.SIC, sub.Name, sub.Tickers[0], sub.Exchanges

	facts := factsStruct.Facts.USGAAP // done to make retrieving data from Facts struct less verbose

	assets := getAssets(facts)

	return nil
}

func getAssets(facts *Facts) int64 {
	return 0
}

func getLatestYearlyFilingValue(fact interface{}) (json.Number, error) {
	v := reflect.ValueOf(fact)

	if v.Kind() != reflect.Struct {
		return json.Number(0), fmt.Errorf("Expected fact to be a struct")
	}

	// Access "Units" Field
	unitsField := v.FieldByName("Units")
	if !unitsField.IsValid() || unitsField.Kind() != reflect.Struct {
		return json.Number(0), fmt.Errorf("Expected a units field of struct type")
	}

	var unitField reflect.Value
	if (unitsField.FieldByName("USD")).IsValid() {
		unitField = unitsField.FieldByName("USD")
	} else if (unitsField.FieldByName("Shares")).IsValid() {
		unitField = unitsField.FieldByName("Shares")
	} else if (unitsField.FieldByName("USDOverShares")).IsValid() {
		unitField = unitsField.FieldByName("USDOverShares")
	} else {
		return json.Number(0), fmt.Errorf("Expected a unitField of USD, Shares, or USDOverShares")
	}

	if unitField.Kind() != reflect.Slice {
		return json.Number(0), fmt.Errorf("Expected unitField to be of type slice")
	}

	return json.Number(0), nil
}

func convertJSONNumberToInt64(num json.Number) int64 {

	return 0
}

func convertJSONNumberToFloat64(num json.Number) float64 {
	return 0.0
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
