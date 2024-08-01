package main

import "encoding/json"

// struct for unmarshalling CIK list JSON
type CIKList struct {
	Data [][]any `json:"data"`
}

// struct for unmarshalling Submissions data JSON
type Submissions struct {
	SIC            string   `json:"sic"` // COULD CAUSE ERROR
	SICDescription string   `json:"sicDescription"`
	Name           string   `json:"name"`
	Tickers        []string `json:"tickers"`
	Exchanges      []string `json:"exchanges"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
}

type FactUnit struct {
	Val   json.Number `json:"val"`
	FP    string      `json:"fp"`
	Form  string      `json:"form"`
	Frame string      `json:"frame"`
}

// struct for unmarshalling Facts data JSON
type Facts struct {
	EntityName string `json:"entityName"`
	Facts      struct {
		USGAAP struct {
			NetIncomeLoss struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"NetIncomeLoss"`

			Assets struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"Assets"`

			Liabilities struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"Liabilities"`

			NetCashProvidedByUsedInOperatingActivities struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"NetCashProvidedByUsedInOperatingActivities"`

			LongTermDebt struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"LongTermDebt"`

			AssetsCurrent struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"AssetsCurrent"`

			LiabilitiesCurrent struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"LiabilitiesCurrent"`

			CommonStockSharesIssued struct {
				Units struct {
					Shares []FactUnit `json:"shares"`
				} `json:"units"`
			} `json:"CommonStockSharesIssued"`

			Revenues struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"Revenues"`

			CostOfGoodsSold struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"CostOfGoodsSold"`

			OperatingIncomeLoss struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"OperatingIncomeLoss"`

			StockholdersEquity struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"StockholdersEquity"`

			CashAndCashEquivalentsAtCarryingValue struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"CashAndCashEquivalentsAtCarryingValue"`

			ShortTermInvestments struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"ShortTermInvestments"`

			AccountsReceivableNetCurrent struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"AccountsReceivableNetCurrent"`

			InterestExpense struct {
				Units struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"InterestExpense"`

			EarningsPerShareBasic struct {
				Units struct {
					USDOverShares []FactUnit `json:"USD/shares"`
				} `json:"units"`
			} `json:"EarningsPerShareBasic"`

			WeightedAverageNumberOfSharesOutstandingBasic struct {
				Units struct {
					Shares []FactUnit `json:"shares"`
				} `json:"units"`
			} `json:"WeightedAverageNumberOfSharesOutstandingBasic"`

			// MIGHT NOT EXIST
			CommonStockDividendsPerShareDeclared struct {
				Units struct {
					USDOverShares []FactUnit `json:"USD/shares"`
				} `json:"units"`
			} `json:"CommonStockDividendsPerShareDeclared"`

			GrossProfit struct {
				Label       string `json:"label"`
				Description string `json:"description"`
				Units       struct {
					USD []FactUnit `json:"USD"`
				} `json:"units"`
			} `json:"GrossProfit"`
		} `json:"us-gaap"`
	} `json:"facts"`
}
