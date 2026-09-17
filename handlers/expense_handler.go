package handlers

import (
	"net/http"
	"strconv"

	"github.com/Nishchal-ll/CloudSpend/models"
	"github.com/gin-gonic/gin"
)

// ListExpensesHandler handles GET /api/expenses
func ListExpensesHandler(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")

	expenses, err := models.GetAllExpenses(category, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

// GetExpenseHandler handles GET /api/expenses/:id
func GetExpenseHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID format"})
		return
	}

	expense, err := models.GetExpenseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expense)
}

// CreateExpenseHandler handles POST /api/expenses
func CreateExpenseHandler(c *gin.Context) {
	var req models.Expense
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := models.CreateExpense(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Expense created successfully",
		"expense": req,
	})
}

// UpdateExpenseHandler handles PUT /api/expenses/:id
func UpdateExpenseHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID format"})
		return
	}

	var req models.Expense
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := models.UpdateExpense(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expense updated successfully",
		"expense": req,
	})
}

// DeleteExpenseHandler handles DELETE /api/expenses/:id
func DeleteExpenseHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID format"})
		return
	}

	if err := models.DeleteExpense(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete expense: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully", "id": id})
}
