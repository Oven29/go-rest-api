package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"go-rest-api/internal/api/dto"
	"go-rest-api/internal/db/model"
	"go-rest-api/internal/db/repository"
)

type TeamService interface {
	CreateTeam(ctx context.Context, req dto.CreateTeamRequest) (*dto.Team, error)
	GetTeam(ctx context.Context, teamName string) (*dto.Team, error)
}

type teamService struct {
	db       *gorm.DB
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamService(db *gorm.DB, teamRepo repository.TeamRepository, userRepo repository.UserRepository) TeamService {
	return &teamService{
		db:       db,
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, req dto.CreateTeamRequest) (*dto.Team, error) {
	var result *dto.Team

	err := s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.teamRepo.ExistsByName(ctx, req.TeamName)
		if err != nil {
			return err
		}
		if exists {
			return &ServiceError{
				Code:    dto.ErrorCodeTeamExists,
				Message: "team_name already exists",
			}
		}

		userIDs := make([]uint, 0, len(req.Members))
		for _, member := range req.Members {
			userID, err := strconv.ParseUint(member.UserID, 10, 32)
			if err != nil {
				var id uint
				_, err = fmt.Sscanf(member.UserID, "u%d", &id)
				if err != nil {
					return fmt.Errorf("invalid user_id format: %s", member.UserID)
				}
				userID = uint64(id)
			}

			user, err := s.userRepo.UpsertUser(ctx, uint(userID), member.Username, member.IsActive)
			if err != nil {
				return err
			}
			userIDs = append(userIDs, user.ID)
		}

		team := &model.Team{
			Name: req.TeamName,
		}
		if err := s.teamRepo.Create(ctx, team); err != nil {
			return err
		}

		for _, userID := range userIDs {
			userTeam := &model.UserTeam{
				UserID: userID,
				TeamID: team.ID,
			}
			if err := tx.Create(userTeam).Error; err != nil {
				return err
			}
		}

		teamWithMembers, err := s.teamRepo.GetByNameWithMembers(ctx, req.TeamName)
		if err != nil {
			return err
		}

		result = mapTeamToDTO(teamWithMembers)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*dto.Team, error) {
	team, err := s.teamRepo.GetByNameWithMembers(ctx, teamName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &ServiceError{
				Code:    dto.ErrorCodeNotFound,
				Message: "team not found",
			}
		}
		return nil, err
	}

	return mapTeamToDTO(team), nil
}

func mapTeamToDTO(team *model.Team) *dto.Team {
	members := make([]dto.TeamMember, len(team.Members))
	for i, member := range team.Members {
		members[i] = dto.TeamMember{
			UserID:   fmt.Sprintf("u%d", member.ID),
			Username: member.Name,
			IsActive: member.IsActive,
		}
	}

	return &dto.Team{
		TeamName: team.Name,
		Members:  members,
	}
}
