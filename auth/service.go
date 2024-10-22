package auth

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

var SECRET_KEY = []byte("BWASTARTUP_s3cr3T_k3Y")

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type jwtService struct{}

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	log.Printf("Validating token: %s", encodedToken)

	// Gunakan jwt.MapClaims
	token, err := jwt.ParseWithClaims(encodedToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return SECRET_KEY, nil
	})

	if err != nil {
		log.Printf("Error during token parsing: %v", err)
		return nil, err
	}

	log.Printf("Token parsed successfully: %v", token.Valid)
	return token, nil
}
