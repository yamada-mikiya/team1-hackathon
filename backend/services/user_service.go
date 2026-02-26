package services

import (
	"context"
	"errors"

	"github.com/yamada-mikiya/team1-hackathon/models"
	"github.com/yamada-mikiya/team1-hackathon/repositories"
)

type UserService interface {
	GetUserDetail(ctx context.Context, id int, authenticatedUserID *int, portfolioKey string) (*models.UserDetailResponse, error)
}

type userService struct {
	userRepo    repositories.UserRepository
	articleRepo repositories.ArticleRepository
}

func NewUserService(userRepo repositories.UserRepository, articleRepo repositories.ArticleRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		articleRepo: articleRepo,
	}
}

// GetUserDetail は公開用のユーザープロフィールと記事一覧を取得します
func (s *userService) GetUserDetail(ctx context.Context, id int, authenticatedUserID *int, portfolioKey string) (*models.UserDetailResponse, error) {
	// 1. Repositoryからユーザーを取得
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. nilチェック (リポジトリで GORM依存を消したため、ここで判定)
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 3. 内部記事を表示するかどうかの判定
	// 認証済みまたは正しいportfolio_keyがある場合は内部記事も表示
	showInternalArticles := authenticatedUserID != nil || (portfolioKey != "" && portfolioKey == user.PortfolioKey)

	// 4. ユーザーの記事を取得
	articles, err := s.articleRepo.FindByAuthorID(id, showInternalArticles)
	if err != nil {
		return nil, err
	}

	// 5. 記事をレスポンス形式に変換
	articleResponses := make([]models.ArticleResponse, 0, len(articles))
	for _, article := range articles {
		tags := make([]string, 0)
		for _, tag := range article.Tags {
			tags = append(tags, tag.Name)
		}

		authorResp := models.AuthorResponse{
			ID:          article.Author.ID,
			Name:        article.Author.Name,
			Affiliation: article.Author.Affiliation,
			IconURL:     article.Author.IconURL,
		}

		articleResp := models.ArticleResponse{
			ID:           article.ID,
			Title:        article.Title,
			ArticleType:  article.ArticleType,
			Content:      article.Content,
			Description:  article.Description,
			ExternalURL:  article.ExternalURL,
			ThumbnailURL: article.ThumbnailURL,
			Slug:         article.Slug,
			Department:   article.Department,
			Status:       article.Status,
			Author:       authorResp,
			CreatedAt:    article.CreatedAt,
			UpdatedAt:    article.UpdatedAt,
			Tags:         tags,
		}
		articleResponses = append(articleResponses, articleResp)
	}

	// 6. レスポンスへの変換
	// 注意: ここで Email や PasswordHash を含めないことでセキュリティを守ります
	res := models.UserDetailResponse{
		ID:          user.ID,
		Name:        user.Name,
		Affiliation: user.Affiliation,
		IconURL:     user.IconURL,
		CreatedAt:   user.CreatedAt,
		Articles:    articleResponses,
	}

	// 認証済みで自分自身のプロフィールを見ている場合のみportfolio_keyを返す
	if authenticatedUserID != nil && *authenticatedUserID == user.ID {
		res.PortfolioKey = &user.PortfolioKey
	}

	return &res, nil
}
