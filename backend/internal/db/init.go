package db

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	db, err = gorm.Open(mysql.Open(user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8&parseTime=True&loc=Local")) /*, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // ロガーをInfoモードに設定
	}*/
	if err != nil {
		panic("failed to connect database")
	}

	//db.AutoMigrate(&model.User{}, &model.BaseIdea{}, &model.Card{}, &model.Idea{})
}
