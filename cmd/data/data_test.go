package main

import (
	"fmt"
	"testing"
)

func TestSeedFunctions(t *testing.T) {

	// testing buildSubmissionsURL()
	t.Run("should pass if url is correctly built", func(t *testing.T) {
		var testCik int64 = 123456
		if url := buildSubmissionsURL(testCik); url != "https://data.sec.gov/submissions/CIK0000123456.json" {
			t.Errorf("expected https://data.sec.gov/submissions/CIK0000123456.json, got %s", url)
		}
	})

	// testing getFactsWithCIK()
	t.Run("should pass if facts data is unmarshalled without error", func(t *testing.T) {
		var cik int64 = 796343
		facts, err := getFactsWithCIK(cik)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		length := len(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares)
		fmt.Println(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares[length-1].Val)
		// length := len(facts.Facts.USGAAP.LongTermDebt.Units.USD)
		// fmt.Println(facts.Facts.USGAAP.LongTermDebt.Units.USD[length - 1].Val)
	})

}
