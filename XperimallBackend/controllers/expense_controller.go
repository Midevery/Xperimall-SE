package controllers

import (
	"XperimallBackend/database"
	"XperimallBackend/models"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type ExpenseInput struct {
	Tenant string  `json:"tenant" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type BulkExpenseInput struct {
	Expenses []ExpenseInput `json:"expenses" binding:"required"`
}

type ExpenseResponse struct {
	ID        uint      `json:"id"`
	Tenant    string    `json:"tenant"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupedExpenseResponse struct {
	Date     string            `json:"date"`
	Total    float64           `json:"total"`
	Expenses []ExpenseResponse `json:"expenses"`
}

func CreateExpense(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input BulkExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expenses []models.Expense
	for _, e := range input.Expenses {
		expense := models.Expense{
			Tenant: e.Tenant,
			Amount: e.Amount,
			UserID: userID.(uint),
		}
		expenses = append(expenses, expense)
	}

	result := database.DB.Create(&expenses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expenses created successfully",
		"data":    expenses,
	})
}

func GetUserExpenses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var expenses []models.Expense
	result := database.DB.Where("user_id = ?", userID.(uint)).Find(&expenses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func GetUserExpensesByDate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var expenses []models.Expense
	result := database.DB.Where("user_id = ? AND deleted_at IS NULL", userID.(uint)).Order("created_at DESC").Find(&expenses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	// Group expenses by date
	groupedExpenses := make(map[string]GroupedExpenseResponse)
	for _, expense := range expenses {
		// Format date to "Monday, 02 January 2006"
		date := expense.CreatedAt.Format("Monday, 02 January 2006")

		// Create a key for the map using the date string
		dateKey := date

		if group, exists := groupedExpenses[dateKey]; exists {
			group.Total += expense.Amount
			group.Expenses = append(group.Expenses, ExpenseResponse{
				ID:        expense.ID,
				Tenant:    expense.Tenant,
				Amount:    expense.Amount,
				CreatedAt: expense.CreatedAt,
			})
			groupedExpenses[dateKey] = group
		} else {
			groupedExpenses[dateKey] = GroupedExpenseResponse{
				Date:  date,
				Total: expense.Amount,
				Expenses: []ExpenseResponse{{
					ID:        expense.ID,
					Tenant:    expense.Tenant,
					Amount:    expense.Amount,
					CreatedAt: expense.CreatedAt,
				}},
			}
		}
	}

	// Convert map to slice and sort by date
	var response []GroupedExpenseResponse
	for _, group := range groupedExpenses {
		response = append(response, group)
	}

	// Sort the response by date in descending order
	sort.Slice(response, func(i, j int) bool {
		dateI, _ := time.Parse("Monday, 02 January 2006", response[i].Date)
		dateJ, _ := time.Parse("Monday, 02 January 2006", response[j].Date)
		return dateI.After(dateJ)
	})

	c.JSON(http.StatusOK, response)
}

func GetUserExpensesByDateDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required"})
		return
	}

	// Parse the date string to time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	var expenses []models.Expense
	// Use DATE() function to compare only the date part, ignoring time
	result := database.DB.Where(
		"user_id = ? AND DATE(created_at) = DATE(?) AND deleted_at IS NULL",
		userID.(uint),
		date,
	).Find(&expenses)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	// Calculate total
	var total float64
	for _, expense := range expenses {
		total += expense.Amount
	}

	response := struct {
		Date     string            `json:"date"`
		Total    float64           `json:"total"`
		Expenses []ExpenseResponse `json:"expenses"`
	}{
		Date:     date.Format("Monday, 02 January 2006"),
		Total:    total,
		Expenses: make([]ExpenseResponse, len(expenses)),
	}

	for i, expense := range expenses {
		response.Expenses[i] = ExpenseResponse{
			ID:        expense.ID,
			Tenant:    expense.Tenant,
			Amount:    expense.Amount,
			CreatedAt: expense.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

func DeleteExpensesByDate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required"})
		return
	}

	// Parse the date string to time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Soft delete by updating deleted_at
	result := database.DB.Model(&models.Expense{}).
		Where("user_id = ? AND DATE(created_at) = DATE(?) AND deleted_at IS NULL", userID.(uint), date).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expenses deleted successfully",
	})
}
