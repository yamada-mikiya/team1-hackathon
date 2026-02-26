package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yamada-mikiya/team1-hackathon/models"
	"github.com/yamada-mikiya/team1-hackathon/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	SignUp(ctx context.Context, req models.SignUpRequest) (models.UserResponse, string, error)
	LogIn(ctx context.Context, req models.AuthenticateRequest) (models.UserResponse, string, error)
	ValidateToken(ctx context.Context, tokenString string) (int, error)
	GetUserByID(ctx context.Context, userID int) (models.UserResponse, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	db        *gorm.DB
	secretKey string
}

func NewAuthService(userRepo repositories.UserRepository, db *gorm.DB, secretKey string) AuthService {
	return &authService{
		userRepo:  userRepo,
		db:        db,
		secretKey: secretKey,
	}
}

// SignUp は新しいユーザーを作成し、JWTトークンを返します
func (s *authService) SignUp(ctx context.Context, req models.SignUpRequest) (models.UserResponse, string, error) {
	// メールアドレスの重複チェック
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.UserResponse{}, "", err
	}
	if existingUser != nil {
		return models.UserResponse{}, "", errors.New("このメールアドレスは既に使用されています")
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, "", err
	}

	// 新しいユーザーを作成
	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Affiliation:  req.Affiliation,
	}

	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		return models.UserResponse{}, "", err
	}

	// JWTトークンを生成
	tokenString, err := s.createToken(ctx, *newUser)
	if err != nil {
		return models.UserResponse{}, "", err
	}

	userResponse := models.UserResponse{
		ID:          newUser.ID,
		Name:        newUser.Name,
		Email:       newUser.Email,
		Affiliation: newUser.Affiliation,
		IconURL:     newUser.IconURL,
	}

	return userResponse, tokenString, nil
}

// LogIn は既存ユーザーを認証し、JWTトークンを返します
func (s *authService) LogIn(ctx context.Context, req models.AuthenticateRequest) (models.UserResponse, string, error) {
	// メールアドレスでユーザーを取得
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserResponse{}, "", errors.New("メールアドレスまたはパスワードが正しくありません")
		}
		return models.UserResponse{}, "", err
	}

	// パスワードを検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return models.UserResponse{}, "", errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// JWTトークンを生成
	tokenString, err := s.createToken(ctx, *user)
	if err != nil {
		return models.UserResponse{}, "", err
	}

	userResponse := models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Affiliation: user.Affiliation,
		IconURL:     user.IconURL,
	}

	return userResponse, tokenString, nil
}

// createToken はJWTトークンを作成します
func (s *authService) createToken(ctx context.Context, user models.User) (string, error) {
	claims := &models.JwtCustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // 72時間後に期限切れ
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken はJWTトークンを検証し、ユーザーIDを返します
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*models.JwtCustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("無効なトークンです")
}

// GetUserByID はユーザーIDからユーザー情報を取得します
func (s *authService) GetUserByID(ctx context.Context, userID int) (models.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return models.UserResponse{}, err
	}

	userResponse := models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Affiliation: user.Affiliation,
		IconURL:     user.IconURL,
	}

	return userResponse, nil
}
