package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/mattrmcg/equitalytics-backend/internal/db"
	"github.com/mattrmcg/equitalytics-backend/internal/services/data"
)

func TestSeedFunctions(t *testing.T) {

	// testing buildSubmissionsURL()
	// t.Run("should pass if url is correctly built", func(t *testing.T) {
	// 	var testCik int64 = 796343
	// 	if url := buildSubmissionsURL(testCik); url != "https://data.sec.gov/submissions/CIK0000796343.json" {
	// 		t.Errorf("expected https://data.sec.gov/submissions/CIK0000123456.json, got %s", url)
	// 	}
	// })

	// t.Run("should pass if submissions is correctly retrieved from url", func(t *testing.T) {
	// 	var testCik int64 = 796343
	// 	// url := buildSubmissionsURL(testCik)

	// 	sub, err := getSubmissionsWithCIK(testCik)
	// 	if err != nil {
	// 		t.Errorf("failed to get Submissions without error: %v\n", err)
	// 	}
	// 	if sub == nil {
	// 		t.Errorf("sub is nil\n")
	// 	}
	// })

	// testing getFactsWithCIK()
	// t.Run("should pass if facts data is unmarshalled without error", func(t *testing.T) {
	// 	var cik int64 = 796343
	// 	facts, err := getFactsWithCIK(cik)
	// 	if err != nil {
	// 		t.Errorf("expected no error, got %v", err)
	// 	}

	// 	length := len(facts.Facts.USGAAP.NetIncomeLoss.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.NetIncomeLoss.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.Assets.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.Assets.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.NetCashProvidedByUsedInOperatingActivities.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.NetCashProvidedByUsedInOperatingActivities.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.LongTermDebt.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.LongTermDebt.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.LiabilitiesCurrent.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.LiabilitiesCurrent.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.CommonStockSharesIssued.Units.Shares)
	// 	fmt.Println(facts.Facts.USGAAP.CommonStockSharesIssued.Units.Shares[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.Revenues.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.Revenues.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.CostOfGoodsSold.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.CostOfGoodsSold.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.OperatingIncomeLoss.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.OperatingIncomeLoss.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.StockholdersEquity.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.StockholdersEquity.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.CashAndCashEquivalentsAtCarryingValue.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.CashAndCashEquivalentsAtCarryingValue.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.ShortTermInvestments.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.ShortTermInvestments.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.AccountsReceivableNetCurrent.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.AccountsReceivableNetCurrent.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.InterestExpense.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.InterestExpense.Units.USD[length-1].Val)

	// 	length = len(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares)
	// 	fmt.Println(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares[length-1].Val)

	// 	// Fixed
	// 	length = len(facts.Facts.USGAAP.WeightedAverageNumberOfSharesOutstandingBasic.Units.Shares)
	// 	fmt.Println(facts.Facts.USGAAP.WeightedAverageNumberOfSharesOutstandingBasic.Units.Shares[length-1].Val)
	// 	//fmt.Printf("Length of WeightedAverageNumberOfSharesOutstandingBasic: %v\n", length)

	// 	// NOT ALWAYS PRESENT
	// 	// length = len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares)
	// 	// fmt.Println(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares[length-1].Val)
	// 	// fmt.Printf("Length of CommonStockDividendsPerShareDeclared: %v\n", length)

	// 	if len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares) == 0 {
	// 		fmt.Println("CommonStockDividendsNotPresent")
	// 	} else {
	// 		length = len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares)
	// 		fmt.Println(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares[length-1].Val)
	// 	}

	// 	length = len(facts.Facts.USGAAP.GrossProfit.Units.USD)
	// 	fmt.Println(facts.Facts.USGAAP.GrossProfit.Units.USD[length-1].Val)

	// 	//dividend := facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared
	// 	//empty := struct{}
	// 	// if &dividend != nil {
	// 	// 	length = len(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares)
	// 	// 	fmt.Println(facts.Facts.USGAAP.CommonStockDividendsPerShareDeclared.Units.USDOverShares[length-1].Val)
	// 	// }

	// 	// length := len(facts.Facts.USGAAP.LongTermDebt.Units.USD)
	// 	// fmt.Println(facts.Facts.USGAAP.LongTermDebt.Units.USD[length - 1].Val)
	// })

}

func TestDataRetrievalFunctions(t *testing.T) {

	// cik := int64(796343)

	// sub, _ := getSubmissionsWithCIK(cik)

	// facts, _ := getFactsWithCIK(cik)
	// t.Run("should pass if correct yearly filing value index is retrieved correctly", func(t *testing.T) {

	// assetsIndex, err := getLatestYearlyFilingValueIndex(facts.Facts.USGAAP.Assets)
	// if err != nil {
	// 	t.Errorf("getLatestYearlyFilingValueIndex failed: %v", err)
	// }
	// EPSIndex, err := getLatestYearlyFilingValueIndex(facts.Facts.USGAAP.EarningsPerShareBasic)
	// if err != nil {
	// 	t.Errorf("getLatestYearlyFilingValueIndex failed for EPS: %v", err)
	// }

	// if EPSIndex == 0 {
	// 	t.Errorf("EPS index was returned as 0")
	// }

	// if assetsIndex == 0 {
	// 	t.Errorf("assets index was returned as 0")
	// }

	// backwardsIndex := len(facts.Facts.USGAAP.Assets.Units.USD) - assetsIndex
	// if backwardsIndex != 5 {
	// 	t.Errorf("expected (counting backwards) index %v, got %v", 5, backwardsIndex)
	// }
	// fmt.Printf("assets index: %v\n", backwardsIndex)

	// backwardsEPSIndex := len(facts.Facts.USGAAP.EarningsPerShareBasic.Units.USDOverShares) - EPSIndex
	// if backwardsEPSIndex != 4 {
	// 	t.Errorf("expected (counting backwards) EPSindex %v, got %v", 4, backwardsEPSIndex)
	// }
	// fmt.Printf("EPS index: %v\n", backwardsEPSIndex)

	// })

	// t.Run("should pass if net income is correctly retrieved", func(t *testing.T) {

	// })

	// t.Run("should pass if assets is correctly retrieved", func(t *testing.T) {
	// 	assets, _, _, err := getAssets(facts)
	// 	if err != nil {
	// 		t.Errorf("unable to retrieve assets: %v", err)
	// 	}
	// 	if assets != 29779000000 {
	// 		t.Errorf("did not receive correct assets amount")
	// 	}
	// })

	// t.Run("should pass if liabilities is correctly retrieved", func(t *testing.T) {
	// 	liabilities, err := getLiabilities(facts)
	// 	if err != nil {
	// 		t.Errorf("unable to retrieve liabilities: %v", err)
	// 	}

	// 	if liabilities != 13261000000 {
	// 		t.Errorf("did not receive correct assets amount")
	// 	}
	// })

	// t.Run("should pass if CommonStockSharesIssued is correctly retrieved", func(t *testing.T) {
	// 	commonStockSharesIssued, previousCommonStockSharesIssued, err := getCurrentAndPreviousYearCommonStockSharesIssued(facts)
	// 	if err != nil {
	// 		t.Errorf("unable to retrieve CommonStockSharesIssued\n")
	// 	}

	// 	if commonStockSharesIssued != 601000000 {
	// 		t.Errorf("did not receive correct current CommonStockSharesIssued amount\n")
	// 	}

	// 	if previousCommonStockSharesIssued != 601000000 {
	// 		t.Errorf("did not receive correct previous CommonStockSharesIssued amount\n")
	// 	}
	// })

	// BROKEN
	// t.Run("should pass if fillCompanyInfoStruct() fills and returns a struct without error", func(t *testing.T) {
	// 	companyInfo, err := fillCompanyInfoStruct(cik, sub, facts)
	// 	if err != nil {
	// 		t.Errorf("fillCompanyInfoStruct() did not exeucte without error\n")
	// 	}
	// 	if companyInfo != nil {
	// 		fmt.Printf("%+v\n", companyInfo)
	// 	}
	// })

	// PG CIK: 80424
	// t.Run("should pass if getCostOfGoodsSold() runs without error for Proctor and Gamble", func(t *testing.T) {
	// 	pgFacts, _ := getFactsWithCIK(int64(80424))

	// 	pgCOGS, err := getCostOfGoodsSold(pgFacts)
	// 	if err != nil {
	// 		t.Errorf("getCostOfGoodsSold() didn't work: %v", err)
	// 	}

	// 	fmt.Println(pgCOGS)

	// })

}

func TestDBFunctions(t *testing.T) {
	dbPool, err := db.CreateDBPool("postgres://root:123@127.0.0.1:5432/eql")
	if err != nil {
		t.Error(err)
	}

	defer db.CloseDBPool(dbPool)

	dataService := data.NewDataService(dbPool)

	t.Run("should pass if db connection is opened successfully", func(t *testing.T) {
		err = dbPool.Ping(context.Background())
		if err != nil {
			t.Errorf("unable to ping db pool: %v", err)
		}
	})
	// t.Run("should pass if company info is correctly added to database", func(t *testing.T) {
	// 	cik := int64(796343)
	// 	sub, _ := getSubmissionsWithCIK(cik)
	// 	facts, _ := getFactsWithCIK(cik)

	// 	companyInfo, err := fillCompanyInfoStruct(cik, sub, facts)
	// 	if err != nil {
	// 		t.Error("unable to fill company info struct")
	// 	}

	// 	exists, err := dataService.CheckIfCIKExists(context.Background(), "0000796343")
	// 	if err != nil {
	// 		t.Error(err)
	// 	}

	// 	if exists != true {

	// 		err = dataService.AddCompanyInfoRow(context.Background(), companyInfo)
	// 		if err != nil {
	// 			t.Errorf("unable to add: %v", err)
	// 		}
	// 	} else {
	// 		fmt.Println("cik already exists!!")
	// 	}
	// })

	// 	t.Run("should pass if exists function works as intended", func(t *testing.T) {

	// 		exists, err := dataService.CheckIfCIKExists(context.Background(), "0000796343")
	// 		if err != nil {
	// 			t.Error(err)
	// 		}

	// 		if exists != true {
	// 			t.Error("didn't return correct boolean value")
	// 		}
	// 	})
	// t.Run("should pass if RetrieveCompanyMarketFacts works as intended", func(t *testing.T) {
	// 	companies, err := dataService.RetrieveCompanyMarketData(context.Background())
	// 	if err != nil {
	// 		t.Error("unable to retrieve market data")
	// 	}
	// 	c := companies[0]
	// 	fmt.Printf("%v, %v, %v, %v\n", c.Ticker, c.EarningsPerShare, c.BookValuePerShare, c.RevenuePerShare)
	// })

	// t.Run("should pass if UpdateMarketPriceFacts() works as intended", func(t *testing.T) {
	// 	err := dataService.UpdateMarketPriceFacts(context.Background(), "0000001750", 1, 1, 1, 1)
	// 	if err != nil {
	// 		t.Error("unable to update market data")
	// 	}
	// })

	t.Run("should pass if fetchMarketPrice() works as intended", func(t *testing.T) {
		mp, err := fetchMarketPrice("SWKS")
		if err != nil {
			t.Error(err)
		}

		fmt.Println(mp)
	})

	t.Run("should pass if update company feature works as intended", func(t *testing.T) {
		sub, err := getSubmissionsWithCIK(int64(1750))
		if err != nil {
			t.Error("unable to get submissions struct")
		}

		facts, err := getFactsWithCIK(int64(1750))
		if err != nil {
			t.Error("unable to get facts struct")
		}

		testInfo, _ := fillCompanyInfoStruct(int64(1750), sub, facts)

		err = dataService.UpdateCompanyInfoRow(context.Background(), testInfo)
		if err != nil {
			t.Errorf("unable to update row: %v", err)
		}
	})
}
