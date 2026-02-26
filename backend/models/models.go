package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Article はブログ記事のモデル
type Article struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AuthorID     int       `json:"author_id" gorm:"not null"`
	ArticleType  string    `json:"article_type" gorm:"type:varchar(50);not null"`
	Title        string    `json:"title" gorm:"type:varchar(255);not null"`
	Content      *string   `json:"content" gorm:"type:text"`
	Description  *string   `json:"description" gorm:"type:text"`
	ExternalURL  *string   `json:"external_url" gorm:"type:text"`
	ThumbnailURL *string   `json:"thumbnail_url" gorm:"type:text"`
	Slug         string    `json:"slug" gorm:"type:varchar(255);unique;not null"`
	Department   string    `json:"department" gorm:"type:varchar(50);not null"`
	Status       string    `json:"status" gorm:"type:varchar(50);not null;default:draft"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Author       *User     `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Tags         []Tag     `json:"tags,omitempty" gorm:"many2many:article_tags;"`
}

// Tag はタグのモデル
type Tag struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"type:varchar(100);unique;not null"`
	IsCategory bool      `json:"is_category" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ArticleTag は記事とタグの中間テーブル
type ArticleTag struct {
	ArticleID int       `json:"article_id" gorm:"primaryKey"`
	TagID     int       `json:"tag_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

// User はユーザーのモデル
type User struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	Email        string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Affiliation  *string   `json:"affiliation" gorm:"type:varchar(255)"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	IconURL      *string   `json:"icon_url" gorm:"type:text"`
	PortfolioKey string    `json:"-" gorm:"type:varchar(255);unique;not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// JwtCustomClaims はJWTのカスタムクレーム
type JwtCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
