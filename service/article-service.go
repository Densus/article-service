package service

import (
	"fmt"
	"github.com/densus/article_service/model/dto"
	"github.com/densus/article_service/model/entity"
	"github.com/densus/article_service/repository"
	"github.com/mashingan/smapping"
)

type ArticleService interface {
	Insert(a dto.CreateArticleDTO) entity.Article
	Update(a dto.UpdateArticleDTO) entity.Article
	Delete(articleID uint64, article entity.Article)
	All() []entity.Article
	FindByID(articleID uint64) entity.Article
	IsAllowedToEdit(authorID string, articleID uint64) bool
}

type articleService struct {
	articleRepository repository.ArticleRepository
}

func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &articleService{articleRepository: articleRepo}
}

func (service *articleService) Insert(a dto.CreateArticleDTO) entity.Article {
	//article := entity.Article{}
	//err := smapping.FillStruct(&article, smapping.MapFields(&a))
	//if err != nil {
	//	log.Fatalf("Failed map %v", err)
	//}
	mapped := smapping.MapFields(&a)
	articleToCreate := entity.Article{}
	err := smapping.FillStruct(&articleToCreate, mapped)
	if err != nil {
		panic(err)
	}
	res := service.articleRepository.Insert(articleToCreate)
	return res
}

func (service *articleService) Update(a dto.UpdateArticleDTO) entity.Article {
	mapped := smapping.MapFields(&a)
	articleToUpdate := entity.Article{}
	err := smapping.FillStruct(&articleToUpdate, mapped)
	if err != nil {
		panic(err)
	}

	res := service.articleRepository.Update(articleToUpdate)
	return res
}

func (service *articleService) Delete(articleID uint64, article entity.Article) {
	service.articleRepository.Delete(articleID, article)
}

func (service *articleService) All() []entity.Article {
	return service.articleRepository.All()
}

func (service *articleService) FindByID(articleID uint64) entity.Article {
	return service.articleRepository.FindByID(articleID)
}

func (service *articleService) IsAllowedToEdit(authorID string, articleID uint64) bool {
	a := service.articleRepository.FindByID(articleID)
	_authorID := fmt.Sprintf("%v", a.AuthorID)

	return authorID == _authorID
}