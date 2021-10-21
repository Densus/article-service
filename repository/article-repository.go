package repository

import (
	"github.com/densus/article_service/model/entity"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Insert(article entity.Article) entity.Article
	Update(article entity.Article) entity.Article
	Delete(articleID uint64, article entity.Article)
	All() []entity.Article
	FindByID(articleID uint64) entity.Article
}

type articleRepository struct {
	dbArticleConnection *gorm.DB
}

func NewArticleRepository(dbArticleConn *gorm.DB) ArticleRepository {
	return &articleRepository{dbArticleConnection: dbArticleConn}
}

func (a *articleRepository) Insert(article entity.Article) entity.Article {
	a.dbArticleConnection.Save(&article)
	a.dbArticleConnection.Preload("Author").Find(&article)
	return article
}

func (a *articleRepository) Update(article entity.Article) entity.Article {
	a.dbArticleConnection.Save(&article)
	a.dbArticleConnection.Preload("Author").Find(&article)
	return article
}

func (a *articleRepository) Delete(articleID uint64, article entity.Article) {
	a.dbArticleConnection.Where("id = ?", articleID).Delete(&article)
}

func (a *articleRepository) All() []entity.Article {
	var articles []entity.Article
	a.dbArticleConnection.Preload("Author").Find(&articles)
	return articles
}

func (a *articleRepository) FindByID(articleID uint64) entity.Article {
	var article entity.Article
	a.dbArticleConnection.Preload("Author").Find(&article, articleID)
	return article
}