package models

import (
	"context"
	"time"
)

type UserService interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type InfoService interface {
	GetInfoByCIK(cik string) (*CompanyInfo, error)
	GetInfoByTicker(ctx context.Context, ticker string) (*CompanyInfo, error)
	GetAllTickers(ctx context.Context) ([]string, error)
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	APIKey    string    `json:"-"`       // figure out optimal way of storing api keys
	Allowed   bool      `json:"allowed"` // allow access to API
	CreatedAt time.Time `json:"createdAt"`
}

// still need to figure out structure of CIK info table
type CompanyInfo struct {
	CIK         string
	SIC         string // possible make this a foreign key
	CompanyName string `json:"companyName"`
	Ticker      string
	Exchanges   []string

	// Auxiliary
	Assets                                   int64 `json:"assets"`
	Liabilities                              int64 `json:"liabilities"`
	Revenues                                 int64 `json:"revenues"`
	PreviousYearRevenues                     int64 `json:"previousYearRevenues"`
	CostOfGoodsSold                          int64 `json:"costOfGoodsSold"`
	GrossProfit                              int64 `json:"grossProfit"`
	PreviousYearGrossProfit                  int64 `json:"previousYearGrossProfit"`
	OperatingIncomeLoss                      int64 `json:"operatingIncomeLoss"`
	StockholdersEquity                       int64 `json:"stockholdersEquity"`
	PreviousYearStockholdersEquity           int64 `json:"previousYearStockholdersEquity"`
	CashAndCashEquivalents                   int64 `json:"cashAndCashEquivalents"`                   // For QuickRatio
	ShortTermInvestments                     int64 `json:"shortTermInvestments"`                     // For QuickRatio
	AccountsReceivableNetCurrent             int64 `json:"accountsReceivableNetCurrent"`             // ForQuickRatio
	PreviousYearAccountsReceivableNetCurrent int64 `json:"previousYearAccountsReceivableNetCurrent"` // For Receivables Turnover Ratio
	InterestExpense                          int64 `json:"interestExpense"`                          // For InterestCoverageRatio

	// InventoryNet not present at some companies - may remove and not display InventoryTurnoverRatio
	// InventoryNet                             int64   `json:"inventoryNet"`                             // For InventoryTurnoverRatio
	// PreviousYearInventoryNet                 int64   `json:"previousYearInventoryNet"`                 // For InventoryTurnoverRatio
	WeightedAverageNumberOfSharesOutstanding int64   `json:"weightedAverageNumberOfSharesOutstanding"` // For PriceToBookRatio
	BookValuePerShare                        float64 `json:"bookValuePerShare"`                        // For PriceToBookRatio
	RevenuePerShare                          float64 `json:"revenuePerShare"`                          // For PriceToSalesRatio
	CommonStockDividendsPerShareDeclared     float64 `json:"commonStockDividensPerShareDeclared"`      // For DividendYield (Might not be reported)

	// PIOTOROSKI SCORING

	PiotroskiScore                         int `json:"piotroskiScore"`
	PointsInProfitability                  int `json:"pointsInProfitability"`
	PointsInLeverageLiquiditySourceOfFunds int `json:"pointsInLeverageLiquiditySourceOfFunds"`
	PointsInOperatingEfficiency            int `json:"pointsInOperatingEfficiency"`

	// POINT 1
	NetIncome           int64 `json:"netIncome"`
	IsNetIncomePositive bool  `json:"isNetIncomePositive"`

	// POINT 2
	ReturnOnAssets           float64 `json:"returnOnAssets"` // NetIncome / Assets
	IsReturnOnAssetsPositive bool    `json:"isReturnOnAssetsPositive"`

	// POINT 3
	OperatingCashFlow           int64 `json:"operatingCashFlow"`
	IsOperatingCashFlowPositive bool  `json:"isOperatingCashFlowPositive"`

	// POINT 4
	IsOCFGreaterThanNetIncome bool `json:"isOCFGreaterThanNetIncome"`

	// POINT 5
	LongTermDebt                    int64 `json:"longTermDebt"`
	PreviousYearLongTermDebt        int64 `json:"previousYearLongTermDebt"`
	IsCurrentLTDLessThanPreviousLTD bool  `json:"isCurrentLTDGreaterThanPreviousLTD"`

	// POINT 6
	AssetsCurrent                    int64   `json:"assetsCurrent"`
	LiabilitiesCurrent               int64   `json:"liabilitiesCurrent"`
	PreviousYearAssetsCurrent        int64   `json:"previousYearAssetsCurrent"`
	PreviousYearLiabilitiesCurrent   int64   `json:"previousYearLiabilitiesCurrent"`
	CurrentRatio                     float64 `json:"currentRatio"`
	PreviousYearCurrentRatio         float64 `json:"previousYearCurrentRatio"`
	IsCurrentCRGreaterThanPreviousCR bool    `json:"isCurrentCRGreaterThanPreviousCR"`

	// POINT 7
	CommonStockSharesIssued             int64 `json:"commonStockShareIssued"`
	PreviousYearCommonStockSharesIssued int64 `json:"previousYearCommonStockShareIssued"`
	SharesIssuedInTheLastYear           int64 `json:"sharesIssuedInTheLastYear"`
	NoNewSharesIssued                   bool  `json:"noNewSharesIssued"`

	// POINT 8
	GrossProfitMargin                float64 `json:"grossProfitMargin"`
	PreviousYearGrossProfitMargin    float64 `json:"previousYearGrossProfitMargin"`
	CurrentGPMGreaterThanPreviousGPM bool    `json:"currentGPMGreaterThanPreviousGPM"`

	// POINT 9
	AssetTurnoverRatio                 float64 `json:"assetTurnoverRatio"`
	PreviousYearAssetTurnoverRatio     float64 `json:"previousYearAssetTurnoverRatio"`
	IsCurrentATRGreaterThanPreviousATR bool    `json:"isCurrentATRGreaterThanPreviousATR"`

	// PROFITABILITY
	// Gross Profit Margin - ALREADY DECLARED
	OperatingProfitMargin float64 `json:"operatingProfitMargin"`
	NetProfitMargin       float64 `json:"netProfitMargin"`
	// ReturnOnAssets - ALREADY DECLARED
	ReturnOnEquity float64 `json:"returnOnEquity"`

	// LIQUIDITY
	// CurrentRatio
	QuickRatio float64 `json:"quickRatio"`

	// SOLVENCY
	DebtToEquityRatio     float64 `json:"debtToEquityRatio"`
	DebtToAssetsRatio     float64 `json:"debtToAssetsRatio"`
	InterestCoverageRatio float64 `json:"interestCoverageRatio"`

	// EFFICIENCY
	// AssetTurnoverRatio - ALREADY DECLARED
	// InventoryTurnoverRatio   float64 `json:"inventoryTurnoverRatio"`
	ReceivablesTurnoverRatio float64 `json:"receivablesTurnoverRatio"`

	// VALUATION
	PriceToEarningsRatio float64 `json:"priceToEarningsRatio"`
	PriceToBookRatio     float64 `json:"priceToBookRatio"`
	PriceToSalesRatio    float64 `json:"priceToSalesRatio"`
	DividendYield        float64 `json:"dividendYield"`
	EarningsPerShare     float64 `json:"earningsPerShare"`

	// MARKET - Data that's pulled from live market reports instead of filings
	MarketPricePerShare float64 `json:"marketPricePerShare"`
	MarketDividendYield float64 `json:"markedDividendYield"`

	// Timestamps for when data is updated
	MarketDataUpdatedAt time.Time `json:"marketDataUpdatedAt"` // live market data update timestamp (probably updated once a day)
	FilingDataUpdatedAt time.Time `json:"filingDataUpdatedAt"` // filing data update timestamp
}

type ReceiveDataPayload struct {
}
