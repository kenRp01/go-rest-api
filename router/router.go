package router

import (
	"go-rest-api/controller"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	// エンドポイントの追加
	e := echo.New()
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	t := e.Group("/tasks")
	// エンドポイントにミドルウェアの追加(デコードし、scho.contextにuserというフィールド名をつけて格納)
	t.Use(echojwt.WithConfig(echojwt.Config{
		// 生成時に指定したJWTトークンを選択
		SigningKey: []byte(os.Getenv("SECRET")),
		// どのトークンを確認するか
		TokenLookup: "cookie:token",
	}))
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)
	return e
}