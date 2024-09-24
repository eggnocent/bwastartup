package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/users"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List the campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	var input campaign.GetCampaignsDetailInput

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	input.ID = id

	campaignDetail, err := h.service.GetCampaignsByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign detail", http.StatusOK, "Success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(users.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("failed to create campaign", http.StatusBadGateway, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success to create campaign", http.StatusOK, "Success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignsDetailInput

	// Log untuk ID dari URI
	log.Println("Handler: Parsing campaign ID from URI")
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		log.Println("Handler: Error parsing campaign ID from URI:", err)
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	log.Println("Handler: Campaign ID from URI:", inputID.ID)

	var inputData campaign.CreateCampaignInput

	// Log untuk JSON yang diterima
	log.Println("Handler: Parsing campaign update data from JSON")
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		log.Println("Handler: Error parsing JSON:", err)
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	log.Println("Handler: Update data:", inputData)

	// Ambil currentUser dari JWT
	currentUser := c.MustGet("currentUser").(users.User)
	log.Println("Handler: Current User ID:", currentUser.ID)
	inputData.User = currentUser

	// Panggil service untuk update campaign
	updateCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		log.Println("Handler: Error updating campaign:", err)
		response := helper.APIResponse("failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	log.Println("Handler: Successfully updated campaign ID:", updateCampaign.ID)
	response := helper.APIResponse("Success to update campaign", http.StatusOK, "Success", campaign.FormatCampaign(updateCampaign))
	c.JSON(http.StatusOK, response)
}
