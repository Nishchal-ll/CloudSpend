package handlers

import (
	"net/http"

	"github.com/Nishchal-ll/CloudSpend/models"
	"github.com/gin-gonic/gin"
)

// RenderLanding handles GET /
func RenderLanding(c *gin.Context) {
	isLoggedIn := false
	if cookie, err := c.Cookie(sessionCookieName); err == nil {
		if _, ok := verifySessionToken(cookie); ok {
			isLoggedIn = true
		}
	}

	c.HTML(http.StatusOK, "landing.html", gin.H{
		"title":      "CloudSpend — Cloud Expense & Budget Intelligence",
		"isLoggedIn": isLoggedIn,
	})
}

// RenderDashboard handles GET /dashboard
func RenderDashboard(c *gin.Context) {
	summary, err := models.GetDashboardSummary()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{
			"error": "Failed to load dashboard data: " + err.Error(),
			"page":  "dashboard",
		})
		return
	}

	user, _ := c.Get("currentUser")

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":   "CloudSpend — Dashboard",
		"page":    "dashboard",
		"summary": summary,
		"user":    user,
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
	user, _ := c.Get("currentUser")

	c.HTML(http.StatusOK, "expenses.html", gin.H{
		"title":            "CloudSpend — Expense Records",
		"page":             "expenses",
		"expenses":         expenses,
		"categories":       categories,
		"providers":        providers,
		"selectedCategory": category,
		"searchTerm":       search,
		"expenseCount":     len(expenses),
		"user":             user,
	})
}
