package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransaction(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("failed to get transaction transaction", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get transaction transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaigns transaction", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransaction(c *gin.Context) {
	// Mendapatkan user saat ini dari context
	currentUser := c.MustGet("currentUser").(users.User)
	userID := currentUser.ID
	log.Println("Handler: Fetching transactions for user ID:", userID)

	// Memanggil service untuk mendapatkan transaksi berdasarkan user ID
	transactions, err := h.service.GetTransactionByUserID(userID)
	if err != nil {
		log.Println("Handler: Error fetching transactions:", err)
		response := helper.APIResponse("failed to get user transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Memformat transaksi yang berhasil didapatkan
	log.Println("Handler: Successfully fetched transactions for user ID:", userID)
	formattedTransactions := transaction.FormatUserTransactions(transactions)

	// Log hasil format untuk debugging
	for _, ft := range formattedTransactions {
		log.Printf("Handler: Formatted transaction ID: %d, Campaign: %s, Image URL: %s",
			ft.ID, ft.Campaign.Name, ft.Campaign.ImageUrl)
	}

	// Mengirim response JSON dengan data transaksi yang sudah diformat
	response := helper.APIResponse("User Transaction", http.StatusOK, "success", formattedTransactions)
	c.JSON(http.StatusOK, response)
}
