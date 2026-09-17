package models

import (
	"database/sql"
	"math"

	"github.com/Nishchal-ll/CloudSpend/database"
)

type Budget struct {
	ID           int64   `json:"id"`
	MonthlyLimit float64 `json:"monthly_limit" binding:"required,gt=0"`
	Currency     string  `json:"currency"`
	UpdatedAt    string  `json:"updated_at"`
}

type CategorySummary struct {
	Category   string  `json:"category"`
	Total      float64 `json:"total"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type ProviderSummary struct {
	Provider   string  `json:"provider"`
	Total      float64 `json:"total"`
	Percentage float64 `json:"percentage"`
}

type DashboardSummary struct {
	MonthlyBudget      float64           `json:"monthly_budget"`
	TotalSpent         float64           `json:"total_spent"`
	RemainingBudget    float64           `json:"remaining_budget"`
	PercentUsed        float64           `json:"percent_used"`
	TotalExpensesCount int               `json:"total_expenses_count"`
	CategoryBreakdown  []CategorySummary `json:"category_breakdown"`
	ProviderBreakdown  []ProviderSummary `json:"provider_breakdown"`
	RecentExpenses     []Expense         `json:"recent_expenses"`
	Currency           string            `json:"currency"`
}

// GetBudget retrieves active budget configuration
func GetBudget() (*Budget, error) {
	var b Budget
	row := database.DB.QueryRow("SELECT id, monthly_limit, currency, updated_at FROM budgets ORDER BY id DESC LIMIT 1")
	err := row.Scan(&b.ID, &b.MonthlyLimit, &b.Currency, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		// Fallback default
		return &Budget{MonthlyLimit: 5000.0, Currency: "USD"}, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// UpdateBudget updates monthly budget limit
func UpdateBudget(newLimit float64) error {
	_, err := database.DB.Exec("UPDATE budgets SET monthly_limit = ?, updated_at = CURRENT_TIMESTAMP WHERE id = (SELECT id FROM budgets ORDER BY id DESC LIMIT 1)", newLimit)
	return err
}

// GetDashboardSummary aggregates metrics and distributions
func GetDashboardSummary() (*DashboardSummary, error) {
	budget, err := GetBudget()
	if err != nil {
		return nil, err
	}

	var totalSpent float64
	var count int
	err = database.DB.QueryRow("SELECT COALESCE(SUM(amount), 0), COUNT(*) FROM expenses").Scan(&totalSpent, &count)
	if err != nil {
		return nil, err
	}

	remaining := budget.MonthlyLimit - totalSpent
	var percentUsed float64
	if budget.MonthlyLimit > 0 {
		percentUsed = math.Round((totalSpent/budget.MonthlyLimit)*1000) / 10
	}

	// Category breakdown
	catRows, err := database.DB.Query(`
		SELECT category, COALESCE(SUM(amount), 0) as total, COUNT(*) as cnt
		FROM expenses
		GROUP BY category
		ORDER BY total DESC
	`)
	var categories []CategorySummary
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var cs CategorySummary
			if err := catRows.Scan(&cs.Category, &cs.Total, &cs.Count); err == nil {
				if totalSpent > 0 {
					cs.Percentage = math.Round((cs.Total/totalSpent)*1000) / 10
				}
				categories = append(categories, cs)
			}
		}
	}

	// Provider breakdown (AWS, Azure, GCP, Cloudflare, etc.)
	provRows, err := database.DB.Query(`
		SELECT provider, COALESCE(SUM(amount), 0) as total
		FROM expenses
		GROUP BY provider
		ORDER BY total DESC
	`)
	var providers []ProviderSummary
	if err == nil {
		defer provRows.Close()
		for provRows.Next() {
			var ps ProviderSummary
			if err := provRows.Scan(&ps.Provider, &ps.Total); err == nil {
				if totalSpent > 0 {
					ps.Percentage = math.Round((ps.Total/totalSpent)*1000) / 10
				}
				providers = append(providers, ps)
			}
		}
	}

	// Recent 5 expenses
	recent, err := GetAllExpenses("", "")
	if err == nil && len(recent) > 5 {
		recent = recent[:5]
	}

	summary := &DashboardSummary{
		MonthlyBudget:      budget.MonthlyLimit,
		TotalSpent:         math.Round(totalSpent*100) / 100,
		RemainingBudget:    math.Round(remaining*100) / 100,
		PercentUsed:        percentUsed,
		TotalExpensesCount: count,
		CategoryBreakdown:  categories,
		ProviderBreakdown:  providers,
		RecentExpenses:     recent,
		Currency:           budget.Currency,
	}

	return summary, nil
}
