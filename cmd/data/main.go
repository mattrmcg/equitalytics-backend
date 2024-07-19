package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
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
	// var errors []error                 // for storing non-fatal errors, database will only be accessed if there are no stored errors

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
				_, err := getFactsWithCIK(cikInt)
				if err != nil {
					log.Println(err)
				}

			}

		}
		time.Sleep(250 * time.Millisecond)
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
	time.Sleep(110 * time.Millisecond)
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
	cikStr := strconv.FormatInt(cik, 10)
	for len(cikStr) < 10 {
		cikStr = fmt.Sprint("0", cikStr)
	}

	url := fmt.Sprint("https://data.sec.gov/submissions/CIK", cikStr, ".json")
	return url
}

func getFactsWithCIK(cik int64) (*Facts, error) {
	time.Sleep(110 * time.Millisecond)

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
	cikStr := strconv.FormatInt(cik, 10)
	for len(cikStr) < 10 {
		cikStr = fmt.Sprint("0", cikStr)
	}

	url := fmt.Sprint("https://data.sec.gov/api/xbrl/companyfacts/CIK", cikStr, ".json")
	return url
}

// CURRENTLY NOT USING OPTION 2

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
