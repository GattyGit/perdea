package routes

import (
	"net/http"

	handlers "backend/internal/handlers"

	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	// 認証が必要ないエンドポイント
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	e.GET("/init", handlers.Hello)

	// JWT認証が必要なルートグループを作成
	r := e.Group("/restricted")
	r.Use(jwtMiddleware.WithConfig(jwtMiddleware.Config{
		SigningKey: []byte("my_secret_key"),
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		},
	}))
	// 認証が必要なエンドポイント
	// r.GET("/hello", handlers.Hello)
	r.GET("/getUserInfo/:id", handlers.GetUserInfo)

	// 基本アイデア用のエンドポイント
	rbi := r.Group("/baseIdea")
	rbi.GET("/getBaseIdea/:id", handlers.GetBaseIdea)
	rbi.POST("/createBaseIdea", handlers.CreateBaseIdea)
	rbi.DELETE("/deleteBaseIdea/:base_idea_id", handlers.DeleteBaseIdea)

	// アイデア用のエンドポイント
	ri := r.Group("/idea")
	ri.POST("/generateIdea", handlers.GenerateIdea)
	ri.GET("/getIdea/:id", handlers.GetIdea)
	ri.POST("/toggleCardStatus/:card_id", handlers.ToggleCardStatus)
	ri.DELETE("/deleteCard/:card_id", handlers.DeleteCard)
}
