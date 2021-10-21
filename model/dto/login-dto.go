package dto

//LoginDTO struct is a data model that used by client to POST from /login url
type LoginDTO struct {
	Email string `json:"email" binding:"required" validate:"required,email,lte=80"`
	Password string `json:"password" binding:"required" validate:"required,gte=8,lte=80"`
}
