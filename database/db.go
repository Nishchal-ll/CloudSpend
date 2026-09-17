package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes SQLite database connection and sets up schema & seed data
func InitDB() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "cloudspend.db"
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Optimize SQLite performance & enable foreign keys
	if _, err := DB.Exec(`PRAGMA journal_mode = WAL; PRAGMA foreign_keys = ON;`); err != nil {
		log.Printf("Warning: failed to set pragma flags: %v", err)
	}

	if err := createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	if err := seedDefaultData(); err != nil {
		log.Printf("Notice: seed check completed: %v", err)
	}

	log.Println("SQLite database initialized successfully at", dbPath)
	return DB, nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS budgets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		monthly_limit REAL NOT NULL DEFAULT 5000.0,
		currency TEXT NOT NULL DEFAULT 'USD',
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		category TEXT NOT NULL,
		provider TEXT NOT NULL DEFAULT 'Other',
		amount REAL NOT NULL,
		description TEXT DEFAULT '',
		expense_date TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(schema)
	return err
}

func seedDefaultData() error {
	// Seed budget if empty
	var budgetCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM budgets").Scan(&budgetCount)
	if err != nil {
		return err
	}
	if budgetCount == 0 {
		_, err = DB.Exec("INSERT INTO budgets (monthly_limit, currency) VALUES (5000.0, 'USD')")
		if err != nil {
			return err
		}
	}

	// Seed sample cloud expenses if empty
	var expenseCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM expenses").Scan(&expenseCount)
	if err != nil {
		return err
	}

	if expenseCount == 0 {
		samples := []struct {
			Title       string
			Category    string
			Provider    string
			Amount      float64
			Description string
			Date        string
		}{
			{"Azure App Service (B1 Linux)", "Compute", "Azure", 54.75, "Production web application container hosting", "2026-09-02"},
			{"Azure Container Registry (Basic)", "Containers", "Azure", 5.00, "Docker image registry and storage", "2026-09-05"},
			{"AWS EC2 t3.medium Instances", "Compute", "AWS", 142.50, "Worker node fleet for background tasks", "2026-09-08"},
			{"AWS Aurora PostgreSQL DB", "Database", "AWS", 280.00, "Primary relational database cluster", "2026-09-10"},
			{"Cloudflare Enterprise DNS & WAF", "Networking", "Cloudflare", 20.00, "Global CDN caching and DDoS protection", "2026-09-12"},
			{"GCP Cloud Storage (Nearline)", "Storage", "GCP", 45.30, "Automated off-site database backups", "2026-09-14"},
			{"Datadog Infrastructure Monitoring", "Monitoring", "SaaS", 115.00, "Application metrics & tracing", "2026-09-15"},
			{"GitHub Team Plan + CI/CD Minutes", "DevOps", "GitHub", 44.00, "Version control and automated actions", "2026-09-16"},
		}

		stmt, err := DB.Prepare(`
			INSERT INTO expenses (title, category, provider, amount, description, expense_date)
			VALUES (?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, s := range samples {
			if _, err := stmt.Exec(s.Title, s.Category, s.Provider, s.Amount, s.Description, s.Date); err != nil {
				return err
			}
		}
		log.Println("Seeded initial cloud expense dataset.")
	}

	return nil
}
