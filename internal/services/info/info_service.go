package info

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mattrmcg/equitalytics-backend/internal/models"
)

type InfoService struct {
	db *pgxpool.Pool
}

func NewInfoService(db *pgxpool.Pool) *InfoService {
	return &InfoService{db: db}
}

func (infoService *InfoService) GetInfoByCIK(cik string) (*models.CompanyInfo, error) {
	return nil, nil
}

func (infoService *InfoService) GetInfoByTicker(ctx context.Context, ticker string) (*models.CompanyInfo, error) {
	query := `SELECT * FROM company_info WHERE ticker = $1`

	var info models.CompanyInfo

	// THIS IS TEMPORARY FIX
	// The market_dividend_yield value in the database is NULL because I forgot to initialize it to 0.0 in the update command
	// When Scan tries to scan a NULL value into &info.MarketDividendYield, it throws an error because MarketDividendYield is a float64

	err := infoService.db.QueryRow(ctx, query, ticker).Scan(
		&info.CIK, &info.SIC, &info.CompanyName, &info.Ticker, &info.Exchanges, &info.Assets, &info.Liabilities, &info.Revenues, &info.PreviousYearRevenues, &info.CostOfGoodsSold, &info.GrossProfit, &info.PreviousYearGrossProfit, &info.OperatingIncomeLoss, &info.StockholdersEquity, &info.PreviousYearStockholdersEquity, &info.CashAndCashEquivalents, &info.ShortTermInvestments, &info.AccountsReceivableNetCurrent, &info.PreviousYearAccountsReceivableNetCurrent, &info.InterestExpense, &info.WeightedAverageNumberOfSharesOutstanding, &info.BookValuePerShare, &info.RevenuePerShare, &info.CommonStockDividendsPerShareDeclared, &info.PiotroskiScore, &info.PointsInProfitability, &info.PointsInLeverageLiquiditySourceOfFunds, &info.PointsInOperatingEfficiency, &info.NetIncome, &info.IsNetIncomePositive, &info.ReturnOnAssets, &info.IsReturnOnAssetsPositive, &info.OperatingCashFlow, &info.IsOperatingCashFlowPositive, &info.IsOCFGreaterThanNetIncome, &info.LongTermDebt, &info.PreviousYearLongTermDebt, &info.IsCurrentLTDLessThanPreviousLTD, &info.AssetsCurrent, &info.LiabilitiesCurrent, &info.PreviousYearAssetsCurrent, &info.PreviousYearLiabilitiesCurrent, &info.CurrentRatio, &info.PreviousYearCurrentRatio, &info.IsCurrentCRGreaterThanPreviousCR, &info.CommonStockSharesIssued, &info.PreviousYearCommonStockSharesIssued, &info.SharesIssuedInTheLastYear, &info.NoNewSharesIssued, &info.GrossProfitMargin, &info.PreviousYearGrossProfitMargin, &info.CurrentGPMGreaterThanPreviousGPM, &info.AssetTurnoverRatio, &info.PreviousYearAssetTurnoverRatio, &info.IsCurrentATRGreaterThanPreviousATR, &info.OperatingProfitMargin, &info.NetProfitMargin, &info.ReturnOnEquity, &info.QuickRatio, &info.DebtToEquityRatio, &info.DebtToAssetsRatio, &info.InterestCoverageRatio, &info.ReceivablesTurnoverRatio, &info.PriceToEarningsRatio, &info.PriceToBookRatio, &info.PriceToSalesRatio, &info.DividendYield, &info.EarningsPerShare, &info.MarketPricePerShare, &info.MarketDividendYield, &info.MarketDataUpdatedAt, &info.FilingDataUpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't query row: %v", err)
	}

	return &info, nil
}

func (infoService *InfoService) GetAllTickers(ctx context.Context) ([]string, error) {
	query := `SELECT ticker FROM company_info`

	rows, err := infoService.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var tickers []string

	for rows.Next() {
		var ticker string
		err := rows.Scan(&ticker)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		tickers = append(tickers, ticker)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %v", rows.Err())
	}

	return tickers, nil
}
