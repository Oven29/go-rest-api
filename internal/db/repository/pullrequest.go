package repository

import (
	"context"

	"gorm.io/gorm"

	"go-rest-api/internal/db/model"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *model.PullRequest) error
	GetByID(ctx context.Context, id uint) (*model.PullRequest, error)
	GetByIDWithRelations(ctx context.Context, id uint) (*model.PullRequest, error)
	Update(ctx context.Context, pr *model.PullRequest) error
	Select(ctx context.Context) ([]model.PullRequest, error)
	Delete(ctx context.Context, id uint) error
	ExistsByID(ctx context.Context, id uint) (bool, error)
	GetReviewerPRs(ctx context.Context, reviewerID uint) ([]model.PullRequest, error)

	AddReviewer(ctx context.Context, prID, reviewerID uint) error
	RemoveReviewer(ctx context.Context, prID, reviewerID uint) error
	IsReviewerAssigned(ctx context.Context, prID, reviewerID uint) (bool, error)
}

type pullRequestRepository struct {
	*BaseRepository[model.PullRequest]
	db *gorm.DB
}

func NewPullRequestRepository(db *gorm.DB) PullRequestRepository {
	return &pullRequestRepository{
		BaseRepository: NewBaseRepository[model.PullRequest](db),
		db:             db,
	}
}

func (r *pullRequestRepository) GetByIDWithRelations(ctx context.Context, id uint) (*model.PullRequest, error) {
	var pr model.PullRequest
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Reviewers").
		First(&pr, id).Error
	return &pr, err
}

func (r *pullRequestRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.PullRequest{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

func (r *pullRequestRepository) GetReviewerPRs(ctx context.Context, reviewerID uint) ([]model.PullRequest, error) {
	var prs []model.PullRequest
	err := r.db.WithContext(ctx).
		Joins("JOIN pull_request_reviewer ON pull_request_reviewer.pr_id = pull_requests.id").
		Where("pull_request_reviewer.reviewer_id = ?", reviewerID).
		Preload("Author").
		Find(&prs).Error
	return prs, err
}

func (r *pullRequestRepository) AddReviewer(ctx context.Context, prID, reviewerID uint) error {
	prReviewer := &model.PullRequestReviewer{
		PrID:       prID,
		ReviewerID: reviewerID,
	}
	return r.db.WithContext(ctx).Create(prReviewer).Error
}

func (r *pullRequestRepository) RemoveReviewer(ctx context.Context, prID, reviewerID uint) error {
	return r.db.WithContext(ctx).
		Where("pr_id = ? AND reviewer_id = ?", prID, reviewerID).
		Delete(&model.PullRequestReviewer{}).Error
}

func (r *pullRequestRepository) IsReviewerAssigned(ctx context.Context, prID, reviewerID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.PullRequestReviewer{}).
		Where("pr_id = ? AND reviewer_id = ?", prID, reviewerID).
		Count(&count).Error
	return count > 0, err
}

func (r *pullRequestRepository) WithTx(tx *gorm.DB) *pullRequestRepository {
	return &pullRequestRepository{
		BaseRepository: r.BaseRepository.WithTx(tx),
		db:             tx,
	}
}
