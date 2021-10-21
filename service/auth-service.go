package service

import (
	"github.com/densus/article_service/model/dto"
	"github.com/densus/article_service/model/entity"
	"github.com/densus/article_service/repository"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//AuthService is a contract that what AuthService can do
type AuthService interface {
	CreateUser(user dto.RegisterDTO) entity.User
	VerifyCredential(email, password string) interface{}
	CheckDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

//NewAuthService creates a new instance that represent the AuthService interface
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepository: userRepo}
}

//CreateUser is a method that process data model dto.RegisterDTO to create user
func (a *authService) CreateUser(user dto.RegisterDTO) entity.User {
	mapped := smapping.MapFields(&user)
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, mapped)
	if err != nil {
		panic(err)
	}

	res := a.userRepository.InsertUser(userToCreate)
	return res
}

//VerifyCredential is a method to verify credentials logically
func (a *authService) VerifyCredential(email, password string) interface{} {
	res := a.userRepository.VerifyCredential(email, password)
	log.Println("res: ", res.(entity.User))
	if v, ok := res.(entity.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

//CheckDuplicateEmail is a method to check if the email has been registered before
func (a *authService) CheckDuplicateEmail(email string) bool {
	res := a.userRepository.CheckDuplicateEmail(email)
	return !(res.Error == nil)
}

//comparePassword is a function to compare input password with hashed password from database
func comparePassword(hashedPassword string, plainPassword []byte) bool {
	byteHashedPassword := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHashedPassword, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}