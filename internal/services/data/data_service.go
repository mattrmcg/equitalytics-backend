package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mattrmcg/equitalytics-backend/internal/models"
)

type DataService struct {
	db *pgxpool.Pool
}

func NewDataService(db *pgxpool.Pool) *DataService {
	return &DataService{db: db}
}

func (ds *DataService) AddCompanyInfoRow(ctx context.Context, c *models.CompanyInfo) error {
	// 69 values
	query := `INSERT INTO company_info (
		cik, sic, company_name, ticker, exchanges, assets, liabilities, revenues, previous_year_revenues, cost_of_goods_sold, gross_profit, previous_year_gross_profit, operating_income_loss, stockholders_equity, previous_year_stockholders_equity, cash_and_cash_equivalents, short_term_investments, accounts_receivable_net_current, previous_year_accounts_receivable_net_current, interest_expense, weighted_average_number_of_shares_outstanding, book_value_per_share, revenue_per_share, common_stock_dividends_per_share_declared, piotroski_score, points_in_profitability, points_in_leverage_liquidity_source_of_funds, points_in_operating_efficiency, net_income, is_net_income_positive, return_on_assets, is_return_on_assets_positive, operating_cash_flow, is_operating_cash_flow_positive, is_ocf_greater_than_net_income, long_term_debt, previous_year_long_term_debt, is_current_ltd_less_than_previous_ltd, assets_current, liabilities_current, previous_year_assets_current, previous_year_liabilities_current, current_ratio, previous_year_current_ratio, is_current_cr_greater_than_previous_cr, common_stock_shares_issued, previous_year_common_stock_shares_issued, shares_issued_in_the_last_year, no_new_shares_issued, gross_profit_margin, previous_year_gross_profit_margin, current_gpm_greater_than_previous_gpm, asset_turnover_ratio, previous_year_asset_turnover_ratio, is_current_atr_greater_than_previous_atr, operating_profit_margin, net_profit_margin, return_on_equity, quick_ratio, debt_to_equity_ratio, debt_to_assets_ratio, interest_coverage_ratio, receivables_turnover_ratio, price_to_earnings_ratio, price_to_book_ratio, price_to_sales_ratio, dividend_yield, earnings_per_share, filing_data_updated_at
	) VALUES (
	 	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69
	)`

	_, err := ds.db.Exec(ctx, query, c.CIK, c.SIC, c.CompanyName, c.Ticker, c.Exchanges, c.Assets, c.Liabilities, c.Revenues, c.PreviousYearRevenues, c.CostOfGoodsSold, c.GrossProfit, c.PreviousYearGrossProfit, c.OperatingIncomeLoss, c.StockholdersEquity, c.PreviousYearStockholdersEquity, c.CashAndCashEquivalents, c.ShortTermInvestments, c.AccountsReceivableNetCurrent, c.PreviousYearAccountsReceivableNetCurrent, c.InterestExpense, c.WeightedAverageNumberOfSharesOutstanding, c.BookValuePerShare, c.RevenuePerShare, c.CommonStockDividendsPerShareDeclared, c.PiotroskiScore, c.PointsInProfitability, c.PointsInLeverageLiquiditySourceOfFunds, c.PointsInOperatingEfficiency, c.NetIncome, c.IsNetIncomePositive, c.ReturnOnAssets, c.IsReturnOnAssetsPositive, c.OperatingCashFlow, c.IsOperatingCashFlowPositive, c.IsOCFGreaterThanNetIncome, c.LongTermDebt, c.PreviousYearLongTermDebt, c.IsCurrentLTDLessThanPreviousLTD, c.AssetsCurrent, c.LiabilitiesCurrent, c.PreviousYearAssetsCurrent, c.PreviousYearLiabilitiesCurrent, c.CurrentRatio, c.PreviousYearCurrentRatio, c.IsCurrentCRGreaterThanPreviousCR, c.CommonStockSharesIssued, c.PreviousYearCommonStockSharesIssued, c.SharesIssuedInTheLastYear, c.NoNewSharesIssued, c.GrossProfitMargin, c.PreviousYearGrossProfitMargin, c.CurrentGPMGreaterThanPreviousGPM, c.AssetTurnoverRatio, c.PreviousYearAssetTurnoverRatio, c.IsCurrentATRGreaterThanPreviousATR, c.OperatingProfitMargin, c.NetProfitMargin, c.ReturnOnEquity, c.QuickRatio, c.DebtToEquityRatio, c.DebtToAssetsRatio, c.InterestCoverageRatio, c.ReceivablesTurnoverRatio, c.PriceToEarningsRatio, c.PriceToBookRatio, c.PriceToSalesRatio, c.DividendYield, c.EarningsPerShare, c.FilingDataUpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DataService) UpdateCompanyInfoRow(ctx context.Context, c *models.CompanyInfo) error {
	query := `
		UPDATE company_info
		SET
			assets = $1,
			liabilities = $2,
			revenues = $3,
			previous_year_revenues = $4,
			cost_of_goods_sold = $5,
			gross_profit = $6,
			previous_year_gross_profit = $7,
			operating_income_loss = $8,
			stockholders_equity = $9,
			previous_year_stockholders_equity = $10,
			cash_and_cash_equivalents = $11,
			short_term_investments = $12,
			accounts_receivable_net_current = $13,
			previous_year_accounts_receivable_net_current = $14,
			interest_expense = $15,
			weighted_average_number_of_shares_outstanding = $16,
			book_value_per_share = $17,
			revenue_per_share = $18,
			common_stock_dividends_per_share_declared = $19,
			piotroski_score = $20,
			points_in_profitability = $21,
			points_in_leverage_liquidity_source_of_funds = $22,
			points_in_operating_efficiency = $23,
			net_income = $24,
			is_net_income_positive = $25,
			return_on_assets = $26,
			is_return_on_assets_positive = $27,
			operating_cash_flow = $28,
			is_operating_cash_flow_positive = $29,
			is_ocf_greater_than_net_income = $30,
			long_term_debt = $31,
			previous_year_long_term_debt = $32,
			is_current_ltd_less_than_previous_ltd = $33,
			assets_current = $34,
			liabilities_current = $35,
			previous_year_assets_current = $36,
			previous_year_liabilities_current = $37,
			current_ratio = $38,
			previous_year_current_ratio = $39,
			is_current_cr_greater_than_previous_cr = $40,
			common_stock_shares_issued = $41,
			previous_year_common_stock_shares_issued = $42,
			shares_issued_in_the_last_year = $43,
			no_new_shares_issued = $44,
			gross_profit_margin = $45,
			previous_year_gross_profit_margin = $46,
			current_gpm_greater_than_previous_gpm = $47,
			asset_turnover_ratio = $48,
			previous_year_asset_turnover_ratio = $49,
			is_current_atr_greater_than_previous_atr = $50,
			operating_profit_margin = $51,
			net_profit_margin = $52,
			return_on_equity = $53,
			quick_ratio = $54,
			debt_to_equity_ratio = $55,
			debt_to_assets_ratio = $56,
			interest_coverage_ratio = $57,
			receivables_turnover_ratio = $58,
			earnings_per_share = $59,
			filing_data_updated_at = $60
		WHERE cik = $61;
		`
	_, err := ds.db.Exec(ctx, query, c.Assets, c.Liabilities, c.Revenues, c.PreviousYearRevenues, c.CostOfGoodsSold, c.GrossProfit, c.PreviousYearGrossProfit, c.OperatingIncomeLoss, c.StockholdersEquity, c.PreviousYearStockholdersEquity, c.CashAndCashEquivalents, c.ShortTermInvestments, c.AccountsReceivableNetCurrent, c.PreviousYearAccountsReceivableNetCurrent, c.InterestExpense, c.WeightedAverageNumberOfSharesOutstanding, c.BookValuePerShare, c.RevenuePerShare, c.CommonStockDividendsPerShareDeclared, c.PiotroskiScore, c.PointsInProfitability, c.PointsInLeverageLiquiditySourceOfFunds, c.PointsInOperatingEfficiency, c.NetIncome, c.IsNetIncomePositive, c.ReturnOnAssets, c.IsReturnOnAssetsPositive, c.OperatingCashFlow, c.IsOperatingCashFlowPositive, c.IsOCFGreaterThanNetIncome, c.LongTermDebt, c.PreviousYearLongTermDebt, c.IsCurrentLTDLessThanPreviousLTD, c.AssetsCurrent, c.LiabilitiesCurrent, c.PreviousYearAssetsCurrent, c.PreviousYearLiabilitiesCurrent, c.CurrentRatio, c.PreviousYearCurrentRatio, c.IsCurrentCRGreaterThanPreviousCR, c.CommonStockSharesIssued, c.PreviousYearCommonStockSharesIssued, c.SharesIssuedInTheLastYear, c.NoNewSharesIssued, c.GrossProfitMargin, c.PreviousYearGrossProfitMargin, c.CurrentGPMGreaterThanPreviousGPM, c.AssetTurnoverRatio, c.PreviousYearAssetTurnoverRatio, c.IsCurrentATRGreaterThanPreviousATR, c.OperatingProfitMargin, c.NetProfitMargin, c.ReturnOnEquity, c.QuickRatio, c.DebtToEquityRatio, c.DebtToAssetsRatio, c.InterestCoverageRatio, c.ReceivablesTurnoverRatio, c.EarningsPerShare, c.FilingDataUpdatedAt, c.CIK)
	if err != nil {
		return err
	}
	return nil
}

// Retrieve all cik strings from database
func (ds *DataService) RetrieveAllCIKs(ctx context.Context) ([]string, error) {
	query := `SELECT cik FROM company_info`

	rows, err := ds.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cikSlice []string

	for rows.Next() {
		var cik string

		if err := rows.Scan(&cik); err != nil {
			return nil, fmt.Errorf("failed to scan cik: %v", err)
		}
		cikSlice = append(cikSlice, cik)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %v", rows.Err())
	}

	return cikSlice, nil
}

func (ds *DataService) CheckIfCIKExists(ctx context.Context, cik string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
				SELECT 1
				FROM company_info
				WHERE cik = $1
			  )`

	err := ds.db.QueryRow(ctx, query, cik).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, err
}

func (ds *DataService) RetrieveCompanyMarketData(ctx context.Context) ([]models.CompanyInfo, error) {
	query := `SELECT cik, ticker, earnings_per_share, book_value_per_share, revenue_per_share, common_stock_dividends_per_share_declared FROM company_info`

	rows, err := ds.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve market data: %v", err)
	}
	defer rows.Close()

	var companies []models.CompanyInfo

	for rows.Next() {
		var company models.CompanyInfo
		if err := rows.Scan(&company.CIK, &company.Ticker, &company.EarningsPerShare, &company.BookValuePerShare, &company.RevenuePerShare, &company.CommonStockDividendsPerShareDeclared); err != nil {
			return nil, fmt.Errorf("failed to scan company info: %v", err)
		}
		companies = append(companies, company)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %v", rows.Err())
	}

	return companies, nil
}

func (ds *DataService) UpdateMarketPriceFacts(ctx context.Context, cik string, marketPrice float64, priceToEarningsRatio float64, priceToBookRatio float64, priceToSalesRatio float64) error {
	query := `UPDATE company_info SET market_price_per_share = $1, price_to_earnings_ratio = $2, price_to_book_ratio = $3, price_to_sales_ratio = $4, market_data_updated_at = $5 WHERE cik = $6`
	_, err := ds.db.Exec(ctx, query, marketPrice, priceToEarningsRatio, priceToBookRatio, priceToSalesRatio, time.Now(), cik)
	if err != nil {
		return fmt.Errorf("failed to update market price: %v", err)
	}
	return nil
}
