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

		length := len(facts.Facts.USGAAP.NetIncomeLoss.Units.USD)
		fmt.Println(facts.Facts.USGAAP.NetIncomeLoss.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.Assets.Units.USD)
		fmt.Println(facts.Facts.USGAAP.Assets.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.NetCashProvidedByUsedInOperatingActivities.Units.USD)
		fmt.Println(facts.Facts.USGAAP.NetCashProvidedByUsedInOperatingActivities.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.LongTermDebt.Units.USD)
		fmt.Println(facts.Facts.USGAAP.LongTermDebt.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.LiabilitiesCurrent.Units.USD)
		fmt.Println(facts.Facts.USGAAP.LiabilitiesCurrent.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.CommonStockSharesIssued.Units.Shares)
		fmt.Println(facts.Facts.USGAAP.CommonStockSharesIssued.Units.Shares[length-1].Val)

		length = len(facts.Facts.USGAAP.Revenues.Units.USD)
		fmt.Println(facts.Facts.USGAAP.Revenues.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.CostOfGoodsSold.Units.USD)
		fmt.Println(facts.Facts.USGAAP.CostOfGoodsSold.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.OperatingIncomeLoss.Units.USD)
		fmt.Println(facts.Facts.USGAAP.OperatingIncomeLoss.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.StockholdersEquity.Units.USD)
		fmt.Println(facts.Facts.USGAAP.StockholdersEquity.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.CashAndCashEquivalentsAtCarryingValue.Units.USD)
		fmt.Println(facts.Facts.USGAAP.CashAndCashEquivalentsAtCarryingValue.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.ShortTermInvestments.Units.USD)
		fmt.Println(facts.Facts.USGAAP.ShortTermInvestments.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.AccountsReceivableNetCurrent.Units.USD)
		fmt.Println(facts.Facts.USGAAP.AccountsReceivableNetCurrent.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.InterestExpense.Units.USD)
		fmt.Println(facts.Facts.USGAAP.InterestExpense.Units.USD[length-1].Val)

		length = len(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares)
		fmt.Println(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares[length-1].Val)

		// Fixed
		length = len(facts.Facts.USGAAP.WeightedAverageNumberOfSharesOutstandingBasic.Units.Shares)
		fmt.Println(facts.Facts.USGAAP.WeightedAverageNumberOfSharesOutstandingBasic.Units.Shares[length-1].Val)
		//fmt.Printf("Length of WeightedAverageNumberOfSharesOutstandingBasic: %v\n", length)

		// NOT ALWAYS PRESENT
		// length = len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares)
		// fmt.Println(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares[length-1].Val)
		// fmt.Printf("Length of CommonStockDividendsPerShareDeclared: %v\n", length)

		if len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares) == 0 {
			fmt.Println("CommonStockDividendsNotPresent")
		} else {
			length = len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares)
			fmt.Println(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares[length-1].Val)
		}

		length = len(facts.Facts.USGAAP.GrossProfit.Units.USD)
		fmt.Println(facts.Facts.USGAAP.GrossProfit.Units.USD[length-1].Val)

		//dividend := facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared
		//empty := struct{}
		// if &dividend != nil {
		// 	length = len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares)
		// 	fmt.Println(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares[length-1].Val)
		// }

		// length := len(facts.Facts.USGAAP.LongTermDebt.Units.USD)
		// fmt.Println(facts.Facts.USGAAP.LongTermDebt.Units.USD[length - 1].Val)
	})

}
