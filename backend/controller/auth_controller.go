package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yamada-mikiya/team1-hackathon/config"
	"github.com/yamada-mikiya/team1-hackathon/models"
	"github.com/yamada-mikiya/team1-hackathon/repositories"
	"github.com/yamada-mikiya/team1-hackathon/services"
	"gorm.io/gorm"
)

type AuthController struct {
	service services.AuthService
	config  *config.Config
}

func NewAuthController(cfg *config.Config, db *gorm.DB) *AuthController {
	userRepo := repositories.NewUserRepository(db)
	service := services.NewAuthService(userRepo, db, cfg.SecretKey)
	return &AuthController{
		service: service,
		config:  cfg,
	}
}

// SignUpHandler は新しいユーザーアカウントを作成します
// @Summary      新規ユーザー登録 (Sign Up)
// @Description  新しいユーザーアカウントを作成し、認証トークンとユーザー情報を返します。
// @Tags         認証 (Auth)
// @Accept       json
// @Produce      json
// @Param        payload body models.SignUpRequest true "ユーザー情報 (メールアドレス、パスワード、名前)"
// @Success      201 {object} models.AuthResponse "登録成功。ユーザー情報と認証トークンを返します。"
// @Failure      400 {object} models.ErrorResponse "リクエストボディが不正です"
// @Failure      409 {object} models.ErrorResponse "指定されたメールアドレスは既に使用されています"
// @Failure      500 {object} models.ErrorResponse "サーバー内部でエラーが発生しました"
// @Router       /api/auth/signup [post]
func (c *AuthController) SignUpHandler(ctx echo.Context) error {
	req := models.SignUpRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "リクエストの形式が不正です",
			Message: err.Error(),
		})
	}

	// バリデーション（簡易版）
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "メールアドレス、パスワード、名前は必須です",
		})
	}

	if len(req.Password) < 8 {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "パスワードは8文字以上である必要があります",
		})
	}

	if req.Affiliation != nil {
		affiliation := *req.Affiliation
		if affiliation != "開発" && affiliation != "マーケティング" && affiliation != "組織管理" {
			return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "所属は開発、マーケティング、組織管理のいずれかである必要があります",
			})
		}
	}

	userRes, tokenString, err := c.service.SignUp(ctx.Request().Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "このメールアドレスは既に登録されています" {
			statusCode = http.StatusConflict
		}
		return ctx.JSON(statusCode, models.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
		})
	}

	// CookieにJWTトークンを設定
	isProduction := c.config.Server.Environment == "production"
	sameSite := http.SameSiteLaxMode
	if isProduction {
		// 本番環境ではクロスサイトリクエストのためSameSite=Noneが必要
		sameSite = http.SameSiteNoneMode
	}
	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
		// Domainは設定しない（異なるドメイン間ではクッキー共有不可のため）
		MaxAge:   259200, // 72時間
		HttpOnly: true,
		Secure:   isProduction, // 本番環境ではtrue（HTTPS必須）
		SameSite: sameSite,     // 本番: None、開発: Lax
	}
	ctx.SetCookie(cookie)

	signUpRes := models.AuthResponse{
		Token: tokenString,
		User:  userRes,
	}

	return ctx.JSON(http.StatusCreated, signUpRes)
}

// LogInHandler は既存ユーザーを認証します
// @Summary      ログイン (Log In)
// @Description  既存のユーザーを認証し、新しい認証トークンを発行します。
// @Tags         認証 (Auth)
// @Accept       json
// @Produce      json
// @Param        payload body models.AuthenticateRequest true "ユーザー情報 (メールアドレスとパスワード)"
// @Success      200 {object} models.AuthResponse "認証成功。新しい認証トークンを返します。"
// @Failure      400 {object} models.ErrorResponse "リクエストボディが不正です"
// @Failure      401 {object} models.ErrorResponse "認証に失敗しました (メールアドレスまたはパスワードが正しくありません)"
// @Failure      500 {object} models.ErrorResponse "サーバー内部でエラーが発生しました"
// @Router       /api/auth/login [post]
func (c *AuthController) LogInHandler(ctx echo.Context) error {
	req := models.AuthenticateRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "リクエストの形式が不正です",
			Message: err.Error(),
		})
	}

	// バリデーション（簡易版）
	if req.Email == "" || req.Password == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "メールアドレスとパスワードは必須です",
		})
	}

	userRes, tokenString, err := c.service.LogIn(ctx.Request().Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "メールアドレスまたはパスワードが正しくありません" {
			statusCode = http.StatusUnauthorized
		}
		return ctx.JSON(statusCode, models.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
		})
	}

	// CookieにJWTトークンを設定
	isProduction := c.config.Server.Environment == "production"
	sameSite := http.SameSiteLaxMode
	if isProduction {
		// 本番環境ではクロスサイトリクエストのためSameSite=Noneが必要
		sameSite = http.SameSiteNoneMode
	}
	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
		// Domainは設定しない（異なるドメイン間ではクッキー共有不可のため）
		MaxAge:   259200, // 72時間
		HttpOnly: true,
		Secure:   isProduction, // 本番環境ではtrue（HTTPS必須）
		SameSite: sameSite,     // 本番: None、開発: Lax
	}
	ctx.SetCookie(cookie)

	logInRes := models.AuthResponse{
		Token: tokenString,
		User:  userRes,
	}

	return ctx.JSON(http.StatusOK, logInRes)
}

// GetMeHandler は現在ログインしているユーザー情報を取得します
// @Summary      現在のユーザー情報を取得
// @Description  Cookieからトークンを読み取り、現在ログイン中のユーザー情報を返します。
// @Tags         認証 (Auth)
// @Produce      json
// @Success      200 {object} models.UserResponse "ユーザー情報"
// @Failure      401 {object} models.ErrorResponse "認証されていません"
// @Router       /api/auth/me [get]
func (c *AuthController) GetMeHandler(ctx echo.Context) error {
	// コンテキストからuser_idを取得（ミドルウェアで設定済み）
	userID := ctx.Get("user_id")
	if userID == nil {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "認証されていません",
		})
	}

	// UserIDからユーザー情報を取得
	userIDInt, ok := userID.(int)
	if !ok {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "ユーザーIDの取得に失敗しました",
		})
	}

	userResponse, err := c.service.GetUserByID(ctx.Request().Context(), userIDInt)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "ユーザーが見つかりません",
		})
	}

	return ctx.JSON(http.StatusOK, userResponse)
}

// LogoutHandler はユーザーをログアウトしてクッキーを削除します
// @Summary      ログアウト (Logout)
// @Description  クッキーを削除してログアウトします。
// @Tags         認証 (Auth)
// @Produce      json
// @Success      200 {object} map[string]string "ログアウト成功"
// @Router       /api/auth/logout [post]
func (c *AuthController) LogoutHandler(ctx echo.Context) error {
	isProduction := c.config.Server.Environment == "production"
	sameSite := http.SameSiteLaxMode
	if isProduction {
		sameSite = http.SameSiteNoneMode
	}

	// クッキーを削除（MaxAgeを-1に設定）
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // クッキーを削除
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: sameSite,
	}
	ctx.SetCookie(cookie)

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "ログアウトしました",
	})
}
