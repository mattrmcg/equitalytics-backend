package main

// Two approaches:
// 1) Use JSON Facts API with CIK list JSON
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
