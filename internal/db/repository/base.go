package repository

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func (r *BaseRepository[T]) Update(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *BaseRepository[T]) Select(ctx context.Context) ([]T, error) {
	var list []T
	err := r.db.WithContext(ctx).Find(&list).Error
	return list, err
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

func (r *BaseRepository[T]) WithTx(tx *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: tx}
}
