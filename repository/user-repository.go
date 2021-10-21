package repository

import (
	"github.com/densus/article_service/model/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

//UserRepository is a contract about what UserRepository can do
type UserRepository interface {
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	ProfileUser(UserID string) entity.User
	AllUser() []entity.User
	VerifyCredential(email, password string) interface{}
	CheckDuplicateEmail(email string) (tx *gorm.DB)
}

type dbUserConnection struct {
	dbConnection *gorm.DB
}

//NewUserRepository creates a new instance that represent the UserRepository interface
func NewUserRepository(db *gorm.DB) UserRepository {
	return &dbUserConnection{
		dbConnection: db,
	}
}

//InsertUser is a method to insert/add user data to database
func (d *dbUserConnection) InsertUser(user entity.User) entity.User {
	user.Password = hashPassword([]byte(user.Password))
	d.dbConnection.Save(&user)
	return user
}

//UpdateUser is a method to update user data to database
func (d *dbUserConnection) UpdateUser(user entity.User) entity.User {
	if user.Password != "" {
		user.Password = hashPassword([]byte(user.Password))
	}else {
		//set user.Password with value tempUser.Password obtained form database
		var tempUser entity.User
		d.dbConnection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}
	d.dbConnection.Save(&user)
	return user
}

//ProfileUser is a method to get user profile by UserID
func (d *dbUserConnection) ProfileUser(UserID string) entity.User {
	var user entity.User
	d.dbConnection.Find(&user, UserID)
	return user
}

//AllUser is a method to get all users
func (d *dbUserConnection) AllUser() []entity.User {
	var users []entity.User
	d.dbConnection.Find(&users)
	return users
}

//VerifyCredential is a method to verify user credentials
func (d *dbUserConnection) VerifyCredential(email, password string) interface{} {
	var user entity.User
	res := d.dbConnection.Where("email = ?", email).Take(&user).Error
	if res == nil {
		return user
	}
	return nil
}

//CheckDuplicateEmail is a method to take user by email
func (d *dbUserConnection) CheckDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return d.dbConnection.Where("email = ?", email).Take(&user)
}

//hashPassword is a function to hashing password before saving to database
func hashPassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash password")
	}

	return string(hash)
}

