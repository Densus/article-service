package dto

//UpdateUserDTO is a data model that used by client to PUT(update) data
type UpdateUserDTO struct {
	ID uint64 `json:"id"`
	Name string `json:"name" binding:"required" validate:"required,gte=3,lte=80"`
	Email string `json:"email" binding:"required" validate:"required,email,lte=80"`
	Password string `json:"password" binding:"required" validate:"required,gte=8,lte=80"`
}