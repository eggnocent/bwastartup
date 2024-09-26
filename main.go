package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/users"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Koneksi ke database
	dsn := "egiwira:12345@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// instance
	userRepository := users.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	//transactionRepository := transaction.NewRepository(db)

	userService := users.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)

	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	router := gin.Default()

	router.Static("/images", "./images")
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaign)
	api.GET("/campaigns/:id", campaignHandler.GetCampaigns)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-image", authMiddleware(authService, userService), campaignHandler.UploadImage)
	router.Run()
}

func authMiddleware(authService auth.Service, userService users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Printf("Authorization Header: %s", authHeader)

		if !strings.Contains(authHeader, "Bearer") {
			log.Println("Authorization header does not contain Bearer")
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		} else {
			log.Println("Invalid Authorization format")
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		log.Printf("Extracted Token: %s", tokenString)

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		log.Printf("Token validated successfully: %v", token.Valid)

		// Casting claims ke jwt.MapClaims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Printf("Claims after casting: %+v", claims)

			// Ambil user_id dari claims
			userID, ok := claims["user_id"].(float64)
			if !ok {
				log.Println("Failed to extract user_id from claims")
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}

			log.Printf("User ID from token: %v", userID)

			// Ambil user dari database
			user, err := userService.GetUserByID(int(userID))
			if err != nil {
				log.Printf("User not found: %v", err)
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}

			log.Printf("User found: %v", user)

			// Set user di context
			c.Set("currentUser", user)
		} else {
			log.Println("Failed to access claims or token is invalid")
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
