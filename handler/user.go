package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/users"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService users.Service
	authService auth.Service
}

func NewUserHandler(userService users.Service, authService auth.Service) *UserHandler {
	return &UserHandler{userService, authService}
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

	token, err := h.authService.GenerateToken(newUsers.ID)
	if err != nil {
		response := helper.APIResponse("Register account has failed", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := users.FormatUser(newUsers, token)

	response := helper.APIResponse("account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input users.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("Login account has failed", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := users.FormatUser(loggedinUser, token)
	response := helper.APIResponse("Succesfuly login", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	var input users.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Check email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "server error"}
		response := helper.APIResponse("Check email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{"is_available": isEmailAvailable}

	var metaMessage string
	if isEmailAvailable {
		metaMessage = "Email is available"
	} else {
		metaMessage = "Email has been registered"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID := 3
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("success to upload image", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
