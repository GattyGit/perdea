package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "backend/internal/db"
	service "backend/internal/services"

	"github.com/labstack/echo/v4"
)

func GenerateIdea(c echo.Context) error {

	type InputInfo struct {
		IdeaWord    string `json:"idea_word"`
		SettingWord string `json:"setting_word"`
		WordNum     string `json:"word_num"`
		BaseIdeaID  string `json:"base_idea_id"`
		IdeaID      string `json:"idea_id"`
	}

	var inputInfo InputInfo

	if err := c.Bind(&inputInfo); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to bind input info"})
	}

	baseIdeaID, err := strconv.ParseUint(inputInfo.BaseIdeaID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to parce"})
	}
	baseIdeaIDUint := uint(baseIdeaID)

	ideaID, err := strconv.ParseUint(inputInfo.IdeaID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to parce"})
	}
	ideaIDUint := uint(ideaID)

	inputText := "Give me " + inputInfo.WordNum + ` "` + inputInfo.SettingWord + `"` + " related to " + `"` + inputInfo.IdeaWord + `"` + " in Japanese only, in the form of a JSON array, please. No unnecessary preliminaries, no keys in the output, just values."

	resp, err := service.GenerateIdea(inputText)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate idea"})
	}
	jsonString := resp.Choices[0].Message.Content
	var result []string

	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to marshal json"})
	}

	cardID, err := db.CreateCard(inputInfo.SettingWord, baseIdeaIDUint, ideaIDUint)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create card"})
	}

	if err := db.CreateIdea(result, cardID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create idea"})
	}

	cards, err := db.GetCard(cardID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get card"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"cards": cards,
	})

}

func GetIdea(c echo.Context) error {
	baseIdeaIDStr := c.Param("id")
	baseIdeaID, err := strconv.ParseUint(baseIdeaIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to parce"})
	}

	resp, err := db.GetIdea(baseIdeaID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieved idea"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"base_idea": resp,
	})
}

func ToggleCardStatus(c echo.Context) error {
	type StatusInfo struct {
		Status uint `json:"status"`
	}

	cardIDStr := c.Param("card_id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to parse card_id"})
	}

	var statusInfo StatusInfo

	if err := c.Bind(&statusInfo); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to bind statusInfo"})
	}

	if err := db.ChangeCardStatus(cardID, statusInfo.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to toggle card status"})
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Successfully toggled card status"})
}

func DeleteCard(c echo.Context) error {
	cardIDstr := c.Param("card_id")
	cardID, err := strconv.Atoi(cardIDstr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid id format"})
	}

	if err := db.DeleteCard(uint(cardID)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete card"})
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Successfully delete card"})
}
