package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	// エンドポイントの追加
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// 許可するreactのドメイン、環境変数
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		// 許可するヘッダー一覧
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
		echo.HeaderAccessControlAllowCredentials, echo.HeaderXCSRFToken},
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE"},
		// cookieの送受信可否
		AllowCredentials: true,
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath: "/",
		CookieDomain: os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		// POSTMANで動作確認できるようにするためSameSiteDefaultModeにする（セキュア属性をfalseにするため）
		CookieSameSite: http.SameSiteNoneMode,
		// CookieSameSite: http.SameSiteDefaultMode,
		// CSRFトークンの有効期限（デフォルトは２４時間）
		// CookieMaxAge: 60
	}))
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)
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