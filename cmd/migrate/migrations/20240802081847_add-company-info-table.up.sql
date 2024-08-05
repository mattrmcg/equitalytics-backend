CREATE TABLE IF NOT EXISTS company_info (
    cik VARCHAR(10) PRIMARY KEY,
    sic VARCHAR(10),
    company_name VARCHAR(255) NOT NULL,
    ticker VARCHAR(10),
    exchanges TEXT[],

    assets BIGINT,
    liabilities BIGINT,
    revenues BIGINT,
    previous_year_revenues BIGINT,
    cost_of_goods_sold BIGINT,
    gross_profit BIGINT,
    previous_year_gross_profit BIGINT,
    operating_income_loss BIGINT,
    stockholders_equity BIGINT,
    previous_year_stockholders_equity BIGINT,
    cash_and_cash_equivalents BIGINT,
    short_term_investments BIGINT,
    accounts_receivable_net_current BIGINT,
    previous_year_accounts_receivable_net_current BIGINT,
    interest_expense BIGINT,

    weighted_average_number_of_shares_outstanding BIGINT,
    book_value_per_share DOUBLE PRECISION,
    revenue_per_share DOUBLE PRECISION,
    common_stock_dividends_per_share_declared DOUBLE PRECISION,

    piotroski_score INTEGER,
    points_in_profitability INTEGER,
    points_in_leverage_liquidity_source_of_funds INTEGER,
    points_in_operating_efficiency INTEGER,

    net_income BIGINT,
    is_net_income_positive BOOLEAN,

    return_on_assets DOUBLE PRECISION,
    is_return_on_assets_positive BOOLEAN,

    operating_cash_flow BIGINT,
    is_operating_cash_flow_positive BOOLEAN,

    is_ocf_greater_than_net_income BOOLEAN,

    long_term_debt BIGINT,
    previous_year_long_term_debt BIGINT,
    is_current_ltd_less_than_previous_ltd BOOLEAN,

    assets_current BIGINT,
    liabilities_current BIGINT,
    previous_year_assets_current BIGINT,
    previous_year_liabilities_current BIGINT,
    current_ratio DOUBLE PRECISION,
    previous_year_current_ratio DOUBLE PRECISION,
    is_current_cr_greater_than_previous_cr BOOLEAN,

    common_stock_shares_issued BIGINT,
    previous_year_common_stock_shares_issued BIGINT,
    shares_issued_in_the_last_year BIGINT,
    no_new_shares_issued BOOLEAN,

    gross_profit_margin DOUBLE PRECISION,
    previous_year_gross_profit_margin DOUBLE PRECISION,
    current_gpm_greater_than_previous_gpm BOOLEAN,

    asset_turnover_ratio DOUBLE PRECISION,
    previous_year_asset_turnover_ratio DOUBLE PRECISION,
    is_current_atr_greater_than_previous_atr BOOLEAN,

    operating_profit_margin DOUBLE PRECISION,
    net_profit_margin DOUBLE PRECISION,
    return_on_equity DOUBLE PRECISION,

    quick_ratio DOUBLE PRECISION,

    debt_to_equity_ratio DOUBLE PRECISION,
    debt_to_assets_ratio DOUBLE PRECISION,
    interest_coverage_ratio DOUBLE PRECISION,

    receivables_turnover_ratio DOUBLE PRECISION,

    price_to_earnings_ratio DOUBLE PRECISION,
    price_to_book_ratio DOUBLE PRECISION,
    price_to_sales_ratio DOUBLE PRECISION,
    dividend_yield DOUBLE PRECISION,
    earnings_per_share DOUBLE PRECISION,

    market_price_per_share DOUBLE PRECISION,
    market_dividend_yield DOUBLE PRECISION,

    market_data_updated_at TIMESTAMPTZ,
    filing_data_updated_at TIMESTAMPTZ


);