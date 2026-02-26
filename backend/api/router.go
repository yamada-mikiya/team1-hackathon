package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yamada-mikiya/team1-hackathon/config"
	"github.com/yamada-mikiya/team1-hackathon/controller"
	"github.com/yamada-mikiya/team1-hackathon/repositories"
	"github.com/yamada-mikiya/team1-hackathon/services"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *echo.Echo {
	router := echo.New()

	corsConfig := middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}

	router.Use(middleware.CORSWithConfig(corsConfig))
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// ヘルスチェック
	router.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Swagger UI
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	// コントローラー初期化
	articleController := controller.NewArticleController(db)
	authController := controller.NewAuthController(cfg, db)

	//UserControllerを作成
	userRepo := repositories.NewUserRepository(db)
	articleRepo := repositories.NewArticleRepository(db)
	userService := services.NewUserService(userRepo, articleRepo)
	userController := controller.NewUserController(userService)
	// APIルート
	api := router.Group("/api")
	{
		// 認証関連
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authController.SignUpHandler)
			auth.POST("/login", authController.LogInHandler)
			auth.POST("/logout", authController.LogoutHandler)
			// 認証必須エンドポイント
			auth.GET("/me", authController.GetMeHandler, OptionalAuthMiddleware(cfg.SecretKey))
		}

		// 記事関連（Optional Auth - トークンがあれば認証、なければゲスト扱い）
		articles := api.Group("/articles", OptionalAuthMiddleware(cfg.SecretKey))
		{
			articles.GET("", articleController.GetArticles)
			articles.GET("/:slug", articleController.GetArticleBySlug)
		}

		//ユーザー関連（未認証でもアクセス可能、OptionalAuthで認証状態を取得）
		users := api.Group("/users", OptionalAuthMiddleware(cfg.SecretKey))
		{
			users.GET("/:id", userController.GetUserDetailHandler)
		}
	}

	return router
}
