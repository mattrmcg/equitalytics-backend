package main

import "encoding/json"

// struct for unmarshalling CIK list JSON
type CIKList struct {
	Data [][]any `json:"data"`
}

// struct for unmarshalling Submissions data JSON
type Submissions struct {
	SIC            json.Number `json:"sic"`
	SICDescription string      `json:"sicDescription"`
	Name           string      `json:"name"`
	Tickers        []string    `json:"tickers"`
	Exchanges      []string    `json:"exchanges"`
	Description    string      `json:"description"`
	Category       string      `json:"category"`
}

// struct for unmarshalling Facts data JSON
type Facts struct {
	Facts struct {
		USGAAP struct {
			NetIncomeLoss struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"NetIncomeLoss"`

			Assets struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"Assets"`

			NetCashProvidedByUsedInOperatingActivities struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"NetCashProvidedByUsedInOperatingActivites"`

			LongTermDebt struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"LongTermDebt"`

			AssetsCurrent struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"AssetsCurrent"`

			LiabilitiesCurrent struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"LiabilitiesCurrent"`

			CommonStockSharesIssued struct {
				Units struct {
					Shares []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"shares"`
				} `json:"units"`
			} `json:"CommonStockSharesIssued"`

			Revenues struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"Revenues"`

			CostOfGoodsSold struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"CostOfGoodsSold"`

			OperatingIncomeLoss struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"OperatingIncomeLoss"`

			StockholdersEquity struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"StockholdersEquity"`

			CashAndCashEquivalentsAtCarryingValue struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"CashAndCashEquivalentsAtCarryingValue"`

			ShortTermInvestments struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"ShortTermInvestments"`

			AccountsReceivableNetCurrent struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"AccountsReceivableNetCurrent"`

			InterestExpense struct {
				Units struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"InterestExpense"`

			EarningsPerShareBasic struct {
				Units struct {
					USDOverShares []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD/shares"`
				} `json:"units"`
			} `json:"EarningsPerShareBasic"`

			WeightedAverageNumberOfSharesOutstandingBasic struct {
				Units struct {
					Shares []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"shares"`
			} `json:"WeightedAverageNumberOfSharesOutstandingBasic"`

			CommonStockDividendsPerShareDeclared struct {
				Units struct {
					USDOverShares []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD/shares"`
				} `json:"units"`
			} `json:"CommonStockDividendsPerShareDeclared"`

			GrossProfit struct {
				Label       string `json:"label"`
				Description string `json:"description"`
				Units       struct {
					USD []struct {
						Val   json.Number `json:"val"`
						FP    string      `json:"fp"`
						Frame string      `json:"frame"`
					} `json:"USD"`
				} `json:"units"`
			} `json:"GrossProfit"`
		} `json:"us-gaap"`
	} `json:"facts"`
}
