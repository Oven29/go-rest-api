package repository

import (
	"context"

	"gorm.io/gorm"

	"go-rest-api/internal/db/model"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	GetByID(ctx context.Context, id uint) (*model.Team, error)
	GetByName(ctx context.Context, name string) (*model.Team, error)
	GetByNameWithMembers(ctx context.Context, name string) (*model.Team, error)
	Update(ctx context.Context, team *model.Team) error
	Select(ctx context.Context) ([]model.Team, error)
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type teamRepository struct {
	*BaseRepository[model.Team]
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{
		BaseRepository: NewBaseRepository[model.Team](db),
		db:             db,
	}
}

func (r *teamRepository) GetByName(ctx context.Context, name string) (*model.Team, error) {
	var team model.Team
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&team).Error
	return &team, err
}

func (r *teamRepository) GetByNameWithMembers(ctx context.Context, name string) (*model.Team, error) {
	var team model.Team
	err := r.db.WithContext(ctx).
		Preload("Members").
		Where("name = ?", name).
		First(&team).Error
	return &team, err
}

func (r *teamRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Team{}).
		Where("name = ?", name).
		Count(&count).Error
	return count > 0, err
}

func (r *teamRepository) WithTx(tx *gorm.DB) *teamRepository {
	return &teamRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
		db:             tx,
	}
}
