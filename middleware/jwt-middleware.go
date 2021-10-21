package middleware

import (
	"github.com/densus/article_service/helper"
	"github.com/densus/article_service/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
)

//AuthorizeJWT is a function to authorize the request from client
func AuthorizeJWT(service service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.ErrorResponse("Failed to process request", "token is not found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		//remove Bearer from token
		splitToken := strings.Split(authHeader, " ")
		authHeader = splitToken[1]

		token, err := service.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claim[user_id]: ", claims["user_id"])
			log.Println("Claim[issuer]: ", claims["issuer"])
		}else {
			log.Println(err)
			response := helper.ErrorResponse("token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
