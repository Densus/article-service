package config

import (
	"fmt"
	"github.com/densus/article_service/model/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)


func SetupDBConnection() *gorm.DB {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed connecting to database")
	}

	//Migrate entity model to database
	errMigration := db.AutoMigrate(&entity.Article{}, &entity.User{})
	if errMigration != nil {
		return nil
	}
	return db
}

//CloseDBConnection is a function to close database connection
func CloseDBConnection(db *gorm.DB)  {
	dbMySQL, err := db.DB()
	if err != nil {
		panic(err)
	}

	errCloseDB := dbMySQL.Close()
	if errCloseDB != nil {
		panic("Fail to close connection")
	}
}
