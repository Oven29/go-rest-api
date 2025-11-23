package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"go-rest-api/internal/api/dto"
	"go-rest-api/internal/db/model"
	"go-rest-api/internal/db/repository"
)

type PullRequestService interface {
	CreatePR(ctx context.Context, req dto.CreatePRRequest) (*dto.PullRequest, error)
	MergePR(ctx context.Context, req dto.MergePRRequest) (*dto.PullRequest, error)
	ReassignReviewer(ctx context.Context, req dto.ReassignPRRequest) (*dto.ReassignPRResponse, error)
}

type pullRequestService struct {
	db       *gorm.DB
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPullRequestService(
	db *gorm.DB,
	prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) PullRequestService {
	return &pullRequestService{
		db:       db,
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *pullRequestService) CreatePR(ctx context.Context, req dto.CreatePRRequest) (*dto.PullRequest, error) {
	prID, err := parsePRID(req.PullRequestID)
	if err != nil {
		return nil, err
	}

	authorID, err := parseUserID(req.AuthorID)
	if err != nil {
		return nil, err
	}

	var result *dto.PullRequest

	err = s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.prRepo.ExistsByID(ctx, prID)
		if err != nil {
			return err
		}
		if exists {
			return &ServiceError{
				Code:    dto.ErrorCodePRExists,
				Message: "PR id already exists",
			}
		}

		author, err := s.userRepo.GetByIDWithTeams(ctx, authorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &ServiceError{
					Code:    dto.ErrorCodeNotFound,
					Message: "author not found",
				}
			}
			return err
		}

		if len(author.Teams) == 0 {
			return &ServiceError{
				Code:    dto.ErrorCodeNotFound,
				Message: "author has no team",
			}
		}

		team, err := s.teamRepo.GetByNameWithMembers(ctx, author.Teams[0].Name)
		if err != nil {
			return err
		}

		pr := &model.PullRequest{
			Title:    req.PullRequestName,
			AuthorID: authorID,
			Status:   model.PrStatusOpen,
		}
		if prID != 0 {
			pr.ID = prID
		}
		if err := s.prRepo.Create(ctx, pr); err != nil {
			return err
		}

		reviewers, err := s.selectReviewers(author, team)
		if err != nil {
			return err
		}

		reviewerIDs := make([]string, 0, len(reviewers))
		for _, reviewer := range reviewers {
			if err := s.prRepo.AddReviewer(ctx, pr.ID, reviewer.ID); err != nil {
				return err
			}
			reviewerIDs = append(reviewerIDs, fmt.Sprintf("u%d", reviewer.ID))
		}

		now := time.Now()
		result = &dto.PullRequest{
			PullRequestID:     req.PullRequestID,
			PullRequestName:   req.PullRequestName,
			AuthorID:          req.AuthorID,
			Status:            dto.PullRequestStatusOpen,
			AssignedReviewers: reviewerIDs,
			CreatedAt:         &now,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pullRequestService) MergePR(ctx context.Context, req dto.MergePRRequest) (*dto.PullRequest, error) {
	prID, err := parsePRID(req.PullRequestID)
	if err != nil {
		return nil, err
	}

	var result *dto.PullRequest

	err = s.db.Transaction(func(tx *gorm.DB) error {
		pr, err := s.prRepo.GetByIDWithRelations(ctx, prID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &ServiceError{
					Code:    dto.ErrorCodeNotFound,
					Message: "PR not found",
				}
			}
			return err
		}

		pr.Status = model.PrStatusMerged
		if err := s.prRepo.Update(ctx, pr); err != nil {
			return err
		}

		reviewerIDs := make([]string, len(pr.Reviewers))
		for i, reviewer := range pr.Reviewers {
			reviewerIDs[i] = fmt.Sprintf("u%d", reviewer.ID)
		}

		now := time.Now()
		result = &dto.PullRequest{
			PullRequestID:     req.PullRequestID,
			PullRequestName:   pr.Title,
			AuthorID:          fmt.Sprintf("u%d", pr.AuthorID),
			Status:            dto.PullRequestStatusMerged,
			AssignedReviewers: reviewerIDs,
			MergedAt:          &now,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pullRequestService) ReassignReviewer(ctx context.Context, req dto.ReassignPRRequest) (*dto.ReassignPRResponse, error) {
	prID, err := parsePRID(req.PullRequestID)
	if err != nil {
		return nil, err
	}

	oldUserID, err := parseUserID(req.OldUserID)
	if err != nil {
		return nil, err
	}

	var result *dto.ReassignPRResponse

	err = s.db.Transaction(func(tx *gorm.DB) error {
		pr, err := s.prRepo.GetByIDWithRelations(ctx, prID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &ServiceError{
					Code:    dto.ErrorCodeNotFound,
					Message: "PR not found",
				}
			}
			return err
		}

		if pr.Status == model.PrStatusMerged {
			return &ServiceError{
				Code:    dto.ErrorCodePRMerged,
				Message: "cannot reassign on merged PR",
			}
		}

		isAssigned, err := s.prRepo.IsReviewerAssigned(ctx, prID, oldUserID)
		if err != nil {
			return err
		}
		if !isAssigned {
			return &ServiceError{
				Code:    dto.ErrorCodeNotAssigned,
				Message: "reviewer is not assigned to this PR",
			}
		}

		oldReviewer, err := s.userRepo.GetByIDWithTeams(ctx, oldUserID)
		if err != nil {
			return err
		}

		if len(oldReviewer.Teams) == 0 {
			return &ServiceError{
				Code:    dto.ErrorCodeNoCandidate,
				Message: "old reviewer has no team",
			}
		}

		team, err := s.teamRepo.GetByNameWithMembers(ctx, oldReviewer.Teams[0].Name)
		if err != nil {
			return err
		}

		newReviewer, err := s.findReplacementReviewer(pr, oldReviewer, team)
		if err != nil {
			return err
		}

		if err := s.prRepo.RemoveReviewer(ctx, prID, oldUserID); err != nil {
			return err
		}

		if err := s.prRepo.AddReviewer(ctx, prID, newReviewer.ID); err != nil {
			return err
		}

		pr, err = s.prRepo.GetByIDWithRelations(ctx, prID)
		if err != nil {
			return err
		}

		reviewerIDs := make([]string, len(pr.Reviewers))
		for i, reviewer := range pr.Reviewers {
			reviewerIDs[i] = fmt.Sprintf("u%d", reviewer.ID)
		}

		result = &dto.ReassignPRResponse{
			PR: dto.PullRequest{
				PullRequestID:     req.PullRequestID,
				PullRequestName:   pr.Title,
				AuthorID:          fmt.Sprintf("u%d", pr.AuthorID),
				Status:            dto.PullRequestStatusOpen,
				AssignedReviewers: reviewerIDs,
			},
			ReplacedBy: fmt.Sprintf("u%d", newReviewer.ID),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pullRequestService) selectReviewers(author *model.User, team *model.Team) ([]model.User, error) {
	candidates := make([]model.User, 0)
	for _, member := range team.Members {
		if member.ID != author.ID && member.IsActive {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return []model.User{}, nil
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	maxReviewers := 2
	if len(candidates) < maxReviewers {
		maxReviewers = len(candidates)
	}

	return candidates[:maxReviewers], nil
}

func (s *pullRequestService) findReplacementReviewer(pr *model.PullRequest, oldReviewer *model.User, team *model.Team) (*model.User, error) {
	assignedIDs := make(map[uint]bool)
	for _, reviewer := range pr.Reviewers {
		assignedIDs[reviewer.ID] = true
	}

	candidates := make([]model.User, 0)
	for _, member := range team.Members {
		if member.ID != pr.AuthorID && // не автор
			member.IsActive && // активен
			member.ID != oldReviewer.ID && // не старый ревьювер
			!assignedIDs[member.ID] { // еще не назначен
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return nil, &ServiceError{
			Code:    dto.ErrorCodeNoCandidate,
			Message: "no active replacement candidate in team",
		}
	}

	return &candidates[rand.Intn(len(candidates))], nil
}

func parsePRID(prID string) (uint, error) {
	var id uint
	_, err := fmt.Sscanf(prID, "pr-%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid pull_request_id format: %s", prID)
	}
	return id, nil
}
