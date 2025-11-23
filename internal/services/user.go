package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"go-rest-api/internal/api/dto"
	"go-rest-api/internal/db/repository"
)

type UserService interface {
	SetIsActive(ctx context.Context, req dto.SetIsActiveRequest) (*dto.User, error)
	GetUserReviews(ctx context.Context, userID string) (*dto.GetUserReviewsResponse, error)
}

type userService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository, prRepo repository.PullRequestRepository) UserService {
	return &userService{
		db:       db,
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *userService) SetIsActive(ctx context.Context, req dto.SetIsActiveRequest) (*dto.User, error) {
	userID, err := parseUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	var result *dto.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.userRepo.GetByIDWithTeams(ctx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &ServiceError{
					Code:    dto.ErrorCodeNotFound,
					Message: "user not found",
				}
			}
			return err
		}

		user.IsActive = req.IsActive
		if err := s.userRepo.Update(ctx, user); err != nil {
			return err
		}

		teamName := ""
		if len(user.Teams) > 0 {
			teamName = user.Teams[0].Name
		}

		result = &dto.User{
			UserID:   req.UserID,
			Username: user.Name,
			TeamName: teamName,
			IsActive: user.IsActive,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *userService) GetUserReviews(ctx context.Context, userID string) (*dto.GetUserReviewsResponse, error) {
	uid, err := parseUserID(userID)
	if err != nil {
		return nil, err
	}

	prs, err := s.prRepo.GetReviewerPRs(ctx, uid)
	if err != nil {
		return nil, err
	}

	prList := make([]dto.PullRequestShort, len(prs))
	for i, pr := range prs {
		prList[i] = dto.PullRequestShort{
			PullRequestID:   fmt.Sprintf("pr-%d", pr.ID),
			PullRequestName: pr.Title,
			AuthorID:        fmt.Sprintf("u%d", pr.AuthorID),
			Status:          dto.PullRequestStatus(pr.Status),
		}
	}

	return &dto.GetUserReviewsResponse{
		UserID:       userID,
		PullRequests: prList,
	}, nil
}

func parseUserID(userID string) (uint, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		var id uint
		_, err = fmt.Sscanf(userID, "u%d", &id)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id format: %s", userID)
		}
		return id, nil
	}
	return uint(uid), nil
}
