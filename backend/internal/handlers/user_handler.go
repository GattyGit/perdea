package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "backend/internal/db"
	model "backend/internal/models"
	service "backend/internal/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := service.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	var loginReq LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return err
	}

	email := loginReq.Email
	password := loginReq.Password

	user, err := service.AuthenticateUser(email, password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":    user.ID,
		"token": tokenString,
	})
}

func GetUserInfo(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		fmt.Println("Error converting ID to integer:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid id format"})
	}

	email, err := db.GetEmail(id)
	fmt.Println(email)
	if err != nil {
		fmt.Println("Error retrieving email:", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve base idea"})
	}

	return c.JSON(http.StatusOK, email)
}
