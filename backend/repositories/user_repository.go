package repositories

import (
	"context"
	"errors"

	"github.com/yamada-mikiya/team1-hackathon/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser は新しいユーザーを作成します
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByEmail はメールアドレスからユーザーを取得します
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID はユーザーIDからユーザーを取得します
func (r *userRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("id = ?", userID).Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &user, nil
}
