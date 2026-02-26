package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yamada-mikiya/team1-hackathon/models"
	"github.com/yamada-mikiya/team1-hackathon/services"
)

type UserController interface {
	GetUserDetailHandler(c echo.Context) error
}

type userController struct {
	s services.UserService
}

func NewUserController(s services.UserService) UserController {
	return &userController{s: s}
}

// GetUserDetailHandler はユーザー詳細を取得します
// @Summary      ユーザー詳細取得
// @Description  指定されたIDのユーザー公開プロフィール（名前、所属、アイコン、記事一覧）を取得します。未認証でもアクセス可能ですが、外部公開記事のみ表示されます。
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        portfolio_key query string false "ポートフォリオキー（内部公開記事も表示する場合に必要）"
// @Success      200  {object}  models.UserDetailResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{id} [get]
func (uc *userController) GetUserDetailHandler(c echo.Context) error {
	// 1. URLパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "IDの形式が不正です",
		})
	}

	// 2. 認証状態を確認
	var authenticatedUserID *int
	if userID := c.Get("user_id"); userID != nil {
		uid := userID.(int)
		authenticatedUserID = &uid
	}

	// 3. ポートフォリオキーの取得（クエリパラメータ）
	portfolioKey := c.QueryParam("portfolio_key")

	// 4. Service呼び出し (Contextを渡す)
	res, err := uc.s.GetUserDetail(c.Request().Context(), id, authenticatedUserID, portfolioKey)

	if err != nil {
		// エラー種類に応じて適切なステータスコードを返す
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "ユーザーが見つかりません",
			})
		}

		// その他のエラー (DB接続エラーなど)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "ユーザー情報の取得に失敗しました",
			Message: err.Error(),
		})
	}

	// 5. レスポンスを返す
	return c.JSON(http.StatusOK, res)
}
