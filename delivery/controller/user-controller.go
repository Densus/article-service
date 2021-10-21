package controller

import (
	"fmt"
	"github.com/densus/article_service/helper"
	"github.com/densus/article_service/model/dto"
	"github.com/densus/article_service/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"strings"
)

//UserController is a contract about what UserController can do
type UserController interface {
	Update(ctx *gin.Context)
	Profile(ctx *gin.Context)
	AllUser(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService service.JWTService
}

//NewUserController creates a new instance that represent UserController interface
func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (controller *userController) Update(ctx *gin.Context) {
	var userUpdateDTO dto.UpdateUserDTO
	errDTO := ctx.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		res := helper.ErrorResponse("failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	validate := validator.New()
	//validate input client based on tag validate model dto.UpdateUserDTO
	if errDTO := validate.Struct(userUpdateDTO); errDTO != nil {
		response := helper.ErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	//remove Bearer from token
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, " ")
	authHeader = splitToken[1]

	token, errToken := controller.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	userUpdateDTO.ID = id
	u := controller.userService.Update(userUpdateDTO)
	res := helper.SuccessResponse(true, "OK", u)
	ctx.JSON(http.StatusOK, res)
}

func (controller *userController) Profile(ctx *gin.Context)  {
	//remove Bearer from token
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, " ")
	authHeader = splitToken[1]

	token, err := controller.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	user := controller.userService.Profile(fmt.Sprintf("%v", claims["user_id"]))
	res := helper.SuccessResponse(true, "OK", user)
	ctx.JSON(http.StatusOK, res)
}

func (controller *userController) AllUser(ctx *gin.Context)  {
	users := controller.userService.AllUser()
	res := helper.SuccessResponse(true, "OK", users)
	ctx.JSON(http.StatusOK, res)
}
