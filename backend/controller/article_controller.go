package controller

import (

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yamada-mikiya/team1-hackathon/models"
	"github.com/yamada-mikiya/team1-hackathon/repositories"
	"github.com/yamada-mikiya/team1-hackathon/services"
	"gorm.io/gorm"
)

type ArticleController struct {
	service services.ArticleService
}

func NewArticleController(db *gorm.DB) *ArticleController {
	repo := repositories.NewArticleRepository(db)
	service := services.NewArticleService(repo)
	return &ArticleController{service: service}
}

// @Summary      記事一覧を取得
// @Description  公開されているブログ記事の一覧を取得します。ログイン済みの場合は内部公開記事も含まれます。ページネーション、部署フィルタ、ステータスフィルタをサポートしています。
// @Tags         記事 (Articles)
// @Accept       json
// @Produce      json
// @Param        page query int false "ページ番号 (デフォルト: 1)" default(1)
// @Param        limit query int false "1ページあたりの件数 (デフォルト: 10, 最大: 100)" default(10)
// @Param        department query string false "部署でフィルタ (Dev, MKT, Ops)" Enums(Dev, MKT, Ops)
// @Param        status query string false "ステータスでフィルタ (internal, public, all)。ゲストの場合は常にpublicのみ" Enums(internal, public, all)
// @Success      200 {object} models.ArticleListResponse "記事一覧"
// @Failure      400 {object} models.ErrorResponse "リクエストパラメータが不正です"
// @Failure      500 {object} models.ErrorResponse "サーバー内部でエラーが発生しました"
// @Router       /api/articles [get]
func (ac *ArticleController) GetArticles(c echo.Context) error {
	// ページネーションパラメータを取得
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// ユーザーがログイン済みかチェック
	isAuthenticated := false
	if userID := c.Get("user_id"); userID != nil {
		isAuthenticated = true
	}

	// フィルタパラメータを取得
	filters := repositories.ArticleFilters{
		Department:      c.QueryParam("department"),
		Status:          c.QueryParam("status"),
		IsAuthenticated: isAuthenticated,
	}

	// サービスから記事一覧を取得
	response, err := ac.service.GetArticles(filters, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "記事の取得に失敗しました",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// CreateExternalArticle は外部記事URLからOGP情報を取得して記事を作成します
// @Summary      外部記事を作成
// @Description  外部サイト（Zenn、Noteなど）の記事URLを指定して、OGP情報を取得し記事として登録します。
// @Tags         記事 (Articles)
// @Accept       json
// @Produce      json
// @Param        payload body models.CreateExternalArticleRequest true "外部記事作成リクエスト"
// @Success      201 {object} models.ArticleResponse "作成された記事"
// @Failure      400 {object} models.ErrorResponse "リクエストが不正です"
// @Failure      401 {object} models.ErrorResponse "認証が必要です"
// @Failure      500 {object} models.ErrorResponse "サーバー内部でエラーが発生しました"
// @Router       /api/articles/external [post]
func (ac *ArticleController) CreateExternalArticle(c echo.Context) error {
	// 認証チェック
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "認証が必要です",
		})
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "ユーザーIDの取得に失敗しました",
		})
	}

	// リクエストボディをパース
	var req models.CreateExternalArticleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "リクエストの形式が不正です",
			Message: err.Error(),
		})
	}

	// バリデーション
	if req.URL == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "URLは必須です",
		})
	}

	if req.Department == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "部署は必須です",
		})
	}

	if req.Status == "" {
		req.Status = "public" // デフォルトはpublic
	}

	// 外部記事を作成
	article, err := ac.service.CreateExternalArticle(c.Request().Context(), userIDInt, req.URL, req.Department, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "外部記事の作成に失敗しました",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, article)
}

// GetArticleBySlug はslugを指定して記事を取得します
// @Summary      記事詳細を取得
// @Description  指定されたslugのブログ記事の詳細を取得します。内部公開記事の場合はログインが必要です。
// @Tags         記事 (Articles)
// @Accept       json
// @Produce      json
// @Param        slug path string true "記事のスラグ" example("go-api-development")
// @Success      200 {object} models.ArticleResponse "記事詳細"
// @Failure      403 {object} models.ErrorResponse "内部公開記事にアクセスするにはログインが必要です"
// @Failure      404 {object} models.ErrorResponse "記事が見つかりません"
// @Failure      500 {object} models.ErrorResponse "サーバー内部でエラーが発生しました"
// @Router       /api/articles/{slug} [get]
func (ac *ArticleController) GetArticleBySlug(c echo.Context) error {
	slug := c.Param("slug")

	// ユーザーがログイン済みかチェック
	isAuthenticated := false
	if userID := c.Get("user_id"); userID != nil {
		isAuthenticated = true
	}

	response, err := ac.service.GetArticleBySlug(slug, isAuthenticated)
	if err != nil {
			if err.Error() == "article not found" {
				return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "記事が見つかりません",
			})
		}
		if err.Error() == "内部公開記事にアクセスするにはログインが必要です" {
			return c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "記事の取得に失敗しました",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}