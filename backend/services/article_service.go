package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/yamada-mikiya/team1-hackathon/models"
	"github.com/yamada-mikiya/team1-hackathon/repositories"
)

type ArticleService interface {
	GetArticles(filters repositories.ArticleFilters, page, limit int) (*models.ArticleListResponse, error)
	GetArticleBySlug(slug string, isAuthenticated bool) (*models.ArticleResponse, error)
	CreateExternalArticle(ctx context.Context, userID int, url, department, status string) (*models.ArticleResponse, error)
}


type articleService struct {
	repo       repositories.ArticleRepository
	ogpService OGPService
}

func NewArticleService(repo repositories.ArticleRepository) ArticleService {
	return &articleService{
		repo:       repo,
		ogpService: NewOGPService(),
	}
}

// GetArticles は記事一覧を取得します
func (s *articleService) GetArticles(filters repositories.ArticleFilters, page, limit int) (*models.ArticleListResponse, error) {
	// リポジトリから記事を取得
	filtersInRepository := repositories.ArticleFilters{
		Department: filters.Department,
		Status:     filters.Status,
		IsAuthenticated: filters.IsAuthenticated,
	}
	articles, totalCount, err := s.repo.FindAll(filtersInRepository, page, limit)
	if err != nil {
		return nil, err
	}

	// レスポンスを構築
	articleResponses := make([]models.ArticleResponse, len(articles))
	for i, article := range articles {
		articleResponses[i] = s.convertArticleToResponse(&article)
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	return &models.ArticleListResponse{
		Articles:   articleResponses,
		TotalCount: int(totalCount),
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}


// GetArticleBySlug はslugを指定して記事を取得します
func (s *articleService) GetArticleBySlug(slug string, isAuthenticated bool) (*models.ArticleResponse, error) {
	article, err := s.repo.FindBySlug(slug, isAuthenticated)
	if err != nil {
		return nil, err
	}

	//データがnilだった場合のチェック
	if article == nil {
		return nil, errors.New("article not found")
	}

	res := s.convertArticleToResponse(article)
	return &res, nil
}

// CreateExternalArticle は外部記事URLからOGP情報を取得して記事を作成します
func (s *articleService) CreateExternalArticle(ctx context.Context, userID int, url, department, status string) (*models.ArticleResponse, error) {
	// OGP情報を取得
	ogp, err := s.ogpService.FetchOGP(ctx, url)
	if err != nil {
		return nil, err
	}

	// slugを生成（タイトルから簡易的に生成）
	slug := generateSlugFromURL(url)

	// 記事モデルを作成
	article := &models.Article{
		AuthorID:     userID,
		ArticleType:  "external",
		Title:        ogp.Title,
		Description:  &ogp.Description,
		ExternalURL:  &url,
		ThumbnailURL: &ogp.Image,
		Slug:         slug,
		Department:   department,
		Status:       status,
	}

	// データベースに保存
	if err := s.repo.Create(article); err != nil {
		return nil, err
	}

	// 保存した記事を取得（Authorを含めるため）
	savedArticle, err := s.repo.FindBySlug(article.Slug, true)
	if err != nil {
		return nil, err
	}

	res := s.convertArticleToResponse(savedArticle)
	return &res, nil
}

// 共通の変換ロジック (Helper Method)
// DBモデル(*models.Article)を受け取り、レスポンスモデル(models.ArticleResponse)を返す
func (s *articleService) convertArticleToResponse(article *models.Article) models.ArticleResponse {
	authorResponse := models.AuthorResponse{}

	// Authorのnilチェックと詰め替え
	if article.Author != nil {
		authorResponse = models.AuthorResponse{
			ID:          article.Author.ID,
			Name:        article.Author.Name,
			Affiliation: article.Author.Affiliation,
			IconURL:     article.Author.IconURL,
		}
	}

	return models.ArticleResponse{
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
		Author:       authorResponse,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}
}

// generateSlugFromURL はURLからユニークなslugを生成します
func generateSlugFromURL(url string) string {
	// URLのハッシュを生成（短く、ユニークに）
	hash := sha256.Sum256([]byte(url))
	hashStr := hex.EncodeToString(hash[:])
	// 最初の16文字を使用
	return "external-" + hashStr[:16]
}
