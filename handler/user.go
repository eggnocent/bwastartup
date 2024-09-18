package handler

import (
	"bwastartup/helper"
	"bwastartup/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService users.Service
}

func NewUserHandler(userService users.Service) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	// menangkap input dari user, dan akan di parsing ke service

	var input users.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register account has failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUsers, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Register account has failed", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// token, err := h.jwtService.GenerateToken()
	formatter := users.FormatUser(newUsers, "tokentokentoken")

	response := helper.APIResponse("account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
