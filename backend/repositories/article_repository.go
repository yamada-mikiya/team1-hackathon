package repositories

import (
	"fmt"

	"github.com/yamada-mikiya/team1-hackathon/models"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	FindAll(filters ArticleFilters, page, limit int) ([]models.Article, int64, error)
	FindBySlug(slug string, isAuthenticated bool) (*models.Article, error)
	FindByAuthorID(authorID int, isAuthenticated bool) ([]models.Article, error)
	Create(article *models.Article) error
}

type ArticleFilters struct {
	Department      string
	Status          string
	IsAuthenticated bool
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// FindAll はフィルタとページネーションを適用して記事一覧を取得します
func (r *articleRepository) FindAll(filters ArticleFilters, page, limit int) ([]models.Article, int64, error) {
	var articles []models.Article
	var totalCount int64

	// クエリを構築
	query := r.db.Model(&models.Article{}).Preload("Author")

	// フィルタを適用
	if filters.Department != "" {
		query = query.Where("department = ?", filters.Department)
	}

	// ゲストの場合：強制的にpublicのみ
	if !filters.IsAuthenticated {
		query = query.Where("status = ?", "public")
	} else {
		// メンバー（ログイン済み）の場合：クエリパラメータに応じて制御
		switch filters.Status {
		case "internal":
			// internalのみ
			query = query.Where("status = ?", "internal")
		case "public":
			// publicのみ
			query = query.Where("status = ?", "public")
		case "all", "":
			// 両方（デフォルト）
			query = query.Where("status IN ?", []string{"public", "internal"})
		default:
			// 指定されたステータスで検索（draft等）
			query = query.Where("status = ?", filters.Status)
		}
	}

	// 総件数を取得
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, totalCount, nil
}

// FindBySlug はslugを指定して記事を取得します
func (r *articleRepository) FindBySlug(slug string, isAuthenticated bool) (*models.Article, error) {
	var article models.Article

	// まずは記事を取得（ステータスを問わず）
	if err := r.db.Preload("Author").Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, err
	}

	// 記事のステータスによって閲覧権限をチェック
	switch article.Status {
	case "public":
		// publicは誰でもOK
		return &article, nil
	case "internal":
		// internalはログイン済みのみ
		if !isAuthenticated {
			return nil, fmt.Errorf("内部公開記事にアクセスするにはログインが必要です")
		}
		return &article, nil
	default:
		// draft等のその他のステータスは見つからない扱い（404）
		return nil, gorm.ErrRecordNotFound
	}
}

// FindByAuthorID は指定した著者の記事を取得します
func (r *articleRepository) FindByAuthorID(authorID int, isAuthenticated bool) ([]models.Article, error) {
	var articles []models.Article

	query := r.db.Preload("Author").Preload("Tags").Where("author_id = ?", authorID)

	// ゲストの場合はpublicのみ
	if !isAuthenticated {
		query = query.Where("status = ?", "public")
	} else {
		// ログイン済みの場合はpublicとinternalを取得
		query = query.Where("status IN ?", []string{"public", "internal"})
	}

	if err := query.Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

// Create は新しい記事を作成します
func (r *articleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}
