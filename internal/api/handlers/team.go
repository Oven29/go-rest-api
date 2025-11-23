package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-rest-api/internal/api/dto"
	"go-rest-api/internal/services"
)

type TeamHandler struct {
	teamService services.TeamService
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// CreateTeam POST /team/add
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: err.Error(),
			},
		})
		return
	}

	team, err := h.teamService.CreateTeam(c.Request.Context(), req)
	if err != nil {
		var serviceErr *services.ServiceError
		if errors.As(err, &serviceErr) {
			statusCode := http.StatusBadRequest
			if serviceErr.Code == dto.ErrorCodeNotFound {
				statusCode = http.StatusNotFound
			}
			c.JSON(statusCode, dto.ErrorResponse{
				Error: dto.ErrorDetail{
					Code:    serviceErr.Code,
					Message: serviceErr.Message,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateTeamResponse{
		Team: *team,
	})
}

// GetTeam GET /team/get?team_name=...
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: "team_name query parameter is required",
			},
		})
		return
	}

	team, err := h.teamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		var serviceErr *services.ServiceError
		if errors.As(err, &serviceErr) {
			statusCode := http.StatusNotFound
			c.JSON(statusCode, dto.ErrorResponse{
				Error: dto.ErrorDetail{
					Code:    serviceErr.Code,
					Message: serviceErr.Message,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, *team)
}
