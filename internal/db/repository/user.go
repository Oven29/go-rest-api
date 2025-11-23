package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-rest-api/internal/db/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByIDWithTeams(ctx context.Context, id uint) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Select(ctx context.Context) ([]model.User, error)
	Delete(ctx context.Context, id uint) error
	UpsertUser(ctx context.Context, id uint, name string, isActive bool) (*model.User, error)
}

type userRepository struct {
	*BaseRepository[model.User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[model.User](db),
		db:             db,
	}
}

func (r *userRepository) GetByIDWithTeams(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Preload("Teams").
		First(&user, id).Error
	return &user, err
}

func (r *userRepository) UpsertUser(ctx context.Context, id uint, name string, isActive bool) (*model.User, error) {
	user := &model.User{
		ID:       id,
		Name:     name,
		IsActive: isActive,
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "is_active"}),
		}).
		Create(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) WithTx(tx *gorm.DB) *userRepository {
	return &userRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
		db:             tx,
	}
}
