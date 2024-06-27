package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	db "backend/internal/db"
	model "backend/internal/models"
)

type BaseIdeaRequest struct {
	Title  string `json:"title"`
	UserID string `json:"user_id"` // リクエスト用の文字列型
}

func GetBaseIdea(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Println("Error converting ID to integer:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid id format"})
	}

	idea, err := db.GetBaseIdea(id)
	if err != nil {
		// データベースからの取得に失敗した場合
		fmt.Println("Error retrieving base idea:", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve base idea"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ideas": idea,
	})
}

func CreateBaseIdea(c echo.Context) error {
	req := new(BaseIdeaRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user_id"})
	}

	// BaseIdea構造体に変換
	baseIdea := &model.BaseIdea{
		Title:  req.Title,
		UserID: uint(userID),
	}

	var baseIdeaID uint

	baseIdeaID, err = db.SaveBaseIdea(baseIdea)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create base idea"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"base_idea_id": baseIdeaID,
	})
}

func DeleteBaseIdea(c echo.Context) error {
	idstr := c.Param("base_idea_id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid id format"})
	}

	if err := db.DeleteBaseIdeaRecursively(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete base idea"})
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Successfully delete base idea"})
}
