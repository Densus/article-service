package service

import (
	"github.com/densus/article_service/model/dto"
	"github.com/densus/article_service/model/entity"
	"github.com/densus/article_service/repository"
	"github.com/mashingan/smapping"
)

//UserService is a contract about what UserService can do
type UserService interface {
	Update(user dto.UpdateUserDTO) entity.User
	Profile(UserID string) entity.User
	AllUser() []entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

//NewUserService creates a new instance that represent UserService interface
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepository: userRepo}
}

//Update is a method that update data user with model dto.UpdateUserDTO
func (u *userService) Update(user dto.UpdateUserDTO) entity.User {
	mapped := smapping.MapFields(&user)
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, mapped)
	if err != nil {
		panic(err)
	}

	res := u.userRepository.UpdateUser(userToUpdate)
	return res
}

//Profile is a method to get profile data
func (u *userService) Profile(UserID string) entity.User {
	return u.userRepository.ProfileUser(UserID)
}

//AllUser is a method to get all users
func (u *userService) AllUser() []entity.User {
	return u.userRepository.AllUser()
}

