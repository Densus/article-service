package dto

//UpdateArticleDTO is a data model that used by client when updating article
type UpdateArticleDTO struct {
	ID uint64 `json:"id" binding:"required"`
	Title string `json:"title" binding:"required" validate:"required,gte=5,lte=100"`
	Description string `json:"description" binding:"required" validate:"required,gte=5,lte=100"`
	AuthorID uint64 `json:"author_id,omitempty"` //optional
}

//CreateArticleDTO is a data model that used by client when creating article
type CreateArticleDTO struct {
	Title string `json:"title" binding:"required" validate:"required,gte=5,lte=100"`
	Description string `json:"description" binding:"required" validate:"required,gte=5,lte=100"`
	AuthorID uint64 `json:"author_id,omitempty"` //optional
}
