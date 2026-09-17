package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nishchal-ll/CloudSpend/database"
)

type Expense struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Provider    string  `json:"provider"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	ExpenseDate string  `json:"expense_date" binding:"required"`
	CreatedAt   string  `json:"created_at"`
}

// GetAllExpenses returns expenses with optional search and category filters
func GetAllExpenses(categoryFilter, searchFilter string) ([]Expense, error) {
	query := "SELECT id, title, category, provider, amount, description, expense_date, created_at FROM expenses"
	var conditions []string
	var args []interface{}

	if categoryFilter != "" && categoryFilter != "All" {
		conditions = append(conditions, "category = ?")
		args = append(args, categoryFilter)
	}

	if searchFilter != "" {
		conditions = append(conditions, "(title LIKE ? OR description LIKE ? OR provider LIKE ?)")
		searchTerm := "%" + searchFilter + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY expense_date DESC, id DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]Expense, 0)
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.Title, &e.Category, &e.Provider, &e.Amount, &e.Description, &e.ExpenseDate, &e.CreatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	return expenses, nil
}

// GetExpenseByID retrieves single expense record
func GetExpenseByID(id int64) (*Expense, error) {
	var e Expense
	row := database.DB.QueryRow(`
		SELECT id, title, category, provider, amount, description, expense_date, created_at 
		FROM expenses WHERE id = ?
	`, id)

	err := row.Scan(&e.ID, &e.Title, &e.Category, &e.Provider, &e.Amount, &e.Description, &e.ExpenseDate, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("expense with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// CreateExpense inserts a new expense into SQLite
func CreateExpense(e *Expense) error {
	if e.Provider == "" {
		e.Provider = "Other"
	}
	stmt, err := database.DB.Prepare(`
		INSERT INTO expenses (title, category, provider, amount, description, expense_date)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.Title, e.Category, e.Provider, e.Amount, e.Description, e.ExpenseDate)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

// UpdateExpense modifies an existing expense record
func UpdateExpense(id int64, e *Expense) error {
	if e.Provider == "" {
		e.Provider = "Other"
	}
	stmt, err := database.DB.Prepare(`
		UPDATE expenses 
		SET title = ?, category = ?, provider = ?, amount = ?, description = ?, expense_date = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.Title, e.Category, e.Provider, e.Amount, e.Description, e.ExpenseDate, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("expense with id %d not found", id)
	}
	e.ID = id
	return nil
}

// DeleteExpense removes an expense record
func DeleteExpense(id int64) error {
	res, err := database.DB.Exec("DELETE FROM expenses WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("expense with id %d not found", id)
	}
	return nil
}
