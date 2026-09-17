package handlers

import (
	"net/http"

	"github.com/Nishchal-ll/CloudSpend/models"
	"github.com/gin-gonic/gin"
)

// DashboardStatsHandler handles GET /api/dashboard
func DashboardStatsHandler(c *gin.Context) {
	summary, err := models.GetDashboardSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate dashboard statistics: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

// GetBudgetHandler handles GET /api/budget
func GetBudgetHandler(c *gin.Context) {
	budget, err := models.GetBudget()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve budget: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, budget)
}

// UpdateBudgetHandler handles POST /api/budget
func UpdateBudgetHandler(c *gin.Context) {
	var req struct {
		MonthlyLimit float64 `json:"monthly_limit" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget limit: " + err.Error()})
		return
	}

	if err := models.UpdateBudget(req.MonthlyLimit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget limit: " + err.Error()})
		return
	}

	budget, _ := models.GetBudget()
	c.JSON(http.StatusOK, gin.H{
		"message": "Monthly budget updated successfully",
		"budget":  budget,
	})
}
