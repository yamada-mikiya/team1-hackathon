package models

// AuthenticateRequest はログインリクエスト
type AuthenticateRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
} // @name AuthenticateRequest

// SignUpRequest はサインアップリクエスト
type SignUpRequest struct {
	Email       string  `json:"email" validate:"required,email" example:"user@example.com"`
	Password    string  `json:"password" validate:"required,min=8" example:"password123"`
	Name        string  `json:"name" validate:"required" example:"山田太郎"`
	Affiliation *string `json:"affiliation" validate:"omitempty,oneof=開発 マーケティング 組織管理" example:"開発"`
} // @name SignUpRequest

// CreateExternalArticleRequest は外部記事作成リクエスト
type CreateExternalArticleRequest struct {
	URL        string `json:"url" validate:"required,url" example:"https://zenn.dev/example/articles/sample"`
	Department string `json:"department" validate:"required,oneof=開発 マーケティング 組織管理" example:"開発"`
	Status     string `json:"status" validate:"required,oneof=draft internal public" example:"public"`
} // @name CreateExternalArticleRequest
