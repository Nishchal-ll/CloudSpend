package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Nishchal-ll/CloudSpend/database"
	"github.com/Nishchal-ll/CloudSpend/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize SQLite Database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}
	defer db.Close()

	// Configure Gin environment
	if os.Getenv("GIN_MODE") == "release" || os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Serve static files and HTML templates
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Health check endpoint for Azure App Service & Container Probes
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "CloudSpend",
			"version": "1.0.0",
		})
	})

	// Web UI Routes (Go HTML Templates)
	router.GET("/", handlers.RenderDashboard)
	router.GET("/expenses", handlers.RenderExpenses)

	// REST API Routes
	api := router.Group("/api")
	{
		// Dashboard & Budget
		api.GET("/dashboard", handlers.DashboardStatsHandler)
		api.GET("/budget", handlers.GetBudgetHandler)
		api.POST("/budget", handlers.UpdateBudgetHandler)

		// Expenses CRUD
		api.GET("/expenses", handlers.ListExpensesHandler)
		api.GET("/expenses/:id", handlers.GetExpenseHandler)
		api.POST("/expenses", handlers.CreateExpenseHandler)
		api.PUT("/expenses/:id", handlers.UpdateExpenseHandler)
		api.DELETE("/expenses/:id", handlers.DeleteExpenseHandler)
	}

	// Dynamic Port Binding for Azure App Service ($PORT / $WEBSITES_PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("WEBSITES_PORT")
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 CloudSpend application starting on port :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Fatal: Server failed to start: %v", err)
	}
}
