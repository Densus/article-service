package controller

import (
	"fmt"
	"github.com/densus/article_service/helper"
	"github.com/densus/article_service/model/dto"
	"github.com/densus/article_service/model/entity"
	"github.com/densus/article_service/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"strings"
)

type ArticleController interface {
	All(ctx *gin.Context)
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByID(ctx *gin.Context)
}

type articleController struct {
	articleService service.ArticleService
	jwtService service.JWTService
}

func NewArticleController(articleServ service.ArticleService, jwtServ service.JWTService) ArticleController {
	return &articleController{articleService: articleServ, jwtService: jwtServ}
}

func (a *articleController) All(ctx *gin.Context) {
	var articles []entity.Article = a.articleService.All()
	res := helper.SuccessResponse(true, "OK", articles)
	ctx.JSON(http.StatusOK, res)
}

func (a *articleController) Insert(ctx *gin.Context) {
	var articleCreateDTO dto.CreateArticleDTO
	if errDTO := ctx.ShouldBind(&articleCreateDTO); errDTO != nil {
		res := helper.ErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	validate := validator.New()
	//validate input client based on tag validate model dto.CreateArticleDTO
	if errDTO := validate.Struct(articleCreateDTO); errDTO != nil {
		response := helper.ErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}else {
		authHeader := ctx.GetHeader("Authorization")
		splitToken := strings.Split(authHeader, " ")
		authHeader = splitToken[1]

		token, errToken := a.jwtService.ValidateToken(authHeader)
		if errToken != nil {
			panic(errToken.Error())
		}
		claims := token.Claims.(jwt.MapClaims)
		authorID := fmt.Sprintf("%v", claims["user_id"])

		convertedAuthorID, err := strconv.ParseUint(authorID, 10, 64)
		if err == nil {
			articleCreateDTO.AuthorID = convertedAuthorID
		}
		result := a.articleService.Insert(articleCreateDTO)
		response := helper.SuccessResponse(true, "OK", result)
		ctx.JSON(http.StatusCreated, response)
	}
}

func (a *articleController) Update(ctx *gin.Context) {
	var articleUpdateDTO dto.UpdateArticleDTO
	errDTO := ctx.ShouldBind(&articleUpdateDTO)
	if errDTO != nil {
		res := helper.ErrorResponse("failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	validate := validator.New()
	//validate input client based on tag validate model dto.UpdateArticleDTO
	if errDTO := validate.Struct(articleUpdateDTO); errDTO != nil {
		response := helper.ErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}else {
		authHeader := ctx.GetHeader("Authorization")
		splitToken := strings.Split(authHeader, " ")
		authHeader = splitToken[1]

		token, errToken := a.jwtService.ValidateToken(authHeader)
		if errToken != nil {
			panic(errToken.Error())
		}
		claims := token.Claims.(jwt.MapClaims)
		authorID := fmt.Sprintf("%v", claims["user_id"])

		if a.articleService.IsAllowedToEdit(authorID, articleUpdateDTO.ID) {
			id, errID := strconv.ParseUint(authorID, 10, 64)
			if errID == nil {
				articleUpdateDTO.AuthorID = id
			}
			result := a.articleService.Update(articleUpdateDTO)
			response := helper.SuccessResponse(true, "OK", result)
			ctx.JSON(http.StatusOK, response)
		}else {
			response := helper.ErrorResponse("You don't have permission", "You're not the owner", helper.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		}
	}
}

func (a *articleController) Delete(ctx *gin.Context) {
	var article entity.Article
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.ErrorResponse("Failed to get id", "param id is not found", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}
	articleID := id
	authHeader := ctx.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, " ")
	authHeader = splitToken[1]

	token, errToken := a.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	authorID := fmt.Sprintf("%v", claims["user_id"])
	if a.articleService.IsAllowedToEdit(authorID, articleID) {
		a.articleService.Delete(articleID, article)
		res := helper.SuccessResponse(true, "Deleted", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
	}else {
		response := helper.ErrorResponse("You don't have permission", "You're not the owner", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
	}
}

func (a *articleController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.ErrorResponse("Param id is not found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var article entity.Article = a.articleService.FindByID(id)
	if (article == entity.Article{}) {
		res := helper.ErrorResponse("data not found", "no data with given id", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
	}else {
		res := helper.SuccessResponse(true, "OK!", article)
		ctx.JSON(http.StatusOK, res)
	}
}
