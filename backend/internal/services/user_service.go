package services

import (
	"fmt"

	db "backend/internal/db"
	model "backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return db.SaveUser(user)
}

func AuthenticateUser(email, password string) (*model.User, error) {
	user, err := db.GetUserByEmail(email)
	fmt.Println(user)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}
