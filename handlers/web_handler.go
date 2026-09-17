package handlers

import (
	"net/http"

	"github.com/Nishchal-ll/CloudSpend/models"
	"github.com/gin-gonic/gin"
)

// RenderDashboard handles GET /
func RenderDashboard(c *gin.Context) {
	summary, err := models.GetDashboardSummary()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "Failed to load dashboard data: " + err.Error(),
			"page":  "dashboard",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "CloudSpend — Dashboard",
		"page":    "dashboard",
		"summary": summary,
	})
}

// RenderExpenses handles GET /expenses
func RenderExpenses(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")

	expenses, err := models.GetAllExpenses(category, search)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "expenses.html", gin.H{
			"error": "Failed to load expenses: " + err.Error(),
			"page":  "expenses",
		})
		return
	}

	categories := []string{"All", "Compute", "Database", "Storage", "Networking", "Containers", "Monitoring", "DevOps", "Other"}
	providers := []string{"Azure", "AWS", "GCP", "Cloudflare", "GitHub", "DigitalOcean", "Other"}

	c.HTML(http.StatusOK, "expenses.html", gin.H{
		"title":            "CloudSpend — Cloud Expense Tracker",
		"page":             "expenses",
		"expenses":         expenses,
		"categories":       categories,
		"providers":        providers,
		"selectedCategory": category,
		"searchTerm":       search,
		"expenseCount":     len(expenses),
	})
}
