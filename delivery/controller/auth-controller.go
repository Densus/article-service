package controller

import (
	"github.com/densus/article_service/helper"
	"github.com/densus/article_service/model/dto"
	"github.com/densus/article_service/model/entity"
	"github.com/densus/article_service/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

//AuthController is a contract about what AuthController can do
type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

//this struct used to accommodate all the services needed
type authController struct {
	authService service.AuthService
	jwtService service.JWTService
}

//NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

//Login controller
func (c *authController) Login(ctx *gin.Context)  {
	var loginDTO dto.LoginDTO

	//validate input client
	if errDTO := ctx.ShouldBind(&loginDTO); errDTO != nil {
		response := helper.ErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	validate := validator.New()
	//validate input client based on tag validate model dto.LoginDTO
	if errDTO := validate.Struct(loginDTO); errDTO != nil {
		response := helper.ErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken
		response := helper.SuccessResponse(true, "OK", helper.DataToken{Token: v.Token, Type: "Bearer"})
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.ErrorResponse("please check your credential", "invalid credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

//Register controller
func (c *authController) Register(ctx *gin.Context)  {
	var registerDTO dto.RegisterDTO

	//validate input client
	if errDTO := ctx.ShouldBind(&registerDTO); errDTO != nil {
		response := helper.ErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	validate := validator.New()
	//validate input client based on tag validate model dto.RegisterDTO
	if errDTO := validate.Struct(registerDTO); errDTO != nil {
		response := helper.ErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if !c.authService.CheckDuplicateEmail(registerDTO.Email) { //email already registered
		response := helper.ErrorResponse("failed to process request", "duplicate email", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusConflict, response)
	}else { //email has never been registered
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.SuccessResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusOK, response)
	}
}
