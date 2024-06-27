package db

import (
	model "backend/internal/models"
)

func SaveUser(user *model.User) error {
	return db.Create(user).Error
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

type Email struct {
	Email string
}

func GetEmail(id uint64) (*Email, error) {

	var email Email

	if err := db.Table("users").Model(&model.User{}).Select("email").Where("id = ?", id).First(&email).Error; err != nil {
		return nil, err
	}
	return &email, nil
}
