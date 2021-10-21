package dto

//RegisterDTO is a data model that used by client to POST/ create new user
type RegisterDTO struct {
	Name string `json:"name" binding:"required" validate:"required,gte=3,lte=80"`
	Email string `json:"email" binding:"required" validate:"required,email,lte=80"`
	Password string `json:"password" binding:"required" validate:"required,gte=8,lte=80"`
}
