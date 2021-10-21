package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

//JWTService is a contract about what JWTService can do
type JWTService interface {
	GenerateToken(userID string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer string
}

//NewJWTService method is creates a new instance that represent JWTService interface
func NewJWTService() JWTService {
	return &jwtService{
		secretKey: "agungndp",
		issuer:    getSecretKey(),
	}
}

//getSecretKey is a function to get secret key value from .env file
func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey != "" {
		secretKey = "agungndp"
	}
	return secretKey
}

//GenerateToken is a method to create new token for new session
func (j *jwtService) GenerateToken(UserID string) string {
	claims := &jwtCustomClaim{
		UserID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(), //1 year
			IssuedAt:  time.Now().Unix(),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

//ValidateToken is a method to  validate token
func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token_ *jwt.Token) (interface{}, error) {
		if _, ok := token_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
