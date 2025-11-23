package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-rest-api/internal/api/dto"
	"go-rest-api/internal/services"
)

type PullRequestHandler struct {
	prService services.PullRequestService
}

func NewPullRequestHandler(prService services.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		prService: prService,
	}
}

// CreatePR POST /pullRequest/create
func (h *PullRequestHandler) CreatePR(c *gin.Context) {
	var req dto.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: err.Error(),
			},
		})
		return
	}

	pr, err := h.prService.CreatePR(c.Request.Context(), req)
	if err != nil {
		var serviceErr *services.ServiceError
		if errors.As(err, &serviceErr) {
			statusCode := http.StatusNotFound
			if serviceErr.Code == dto.ErrorCodePRExists {
				statusCode = http.StatusConflict
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

	c.JSON(http.StatusCreated, dto.CreatePRResponse{
		PR: *pr,
	})
}

// MergePR POST /pullRequest/merge
func (h *PullRequestHandler) MergePR(c *gin.Context) {
	var req dto.MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: err.Error(),
			},
		})
		return
	}

	pr, err := h.prService.MergePR(c.Request.Context(), req)
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

	c.JSON(http.StatusOK, dto.MergePRResponse{
		PR: *pr,
	})
}

// ReassignReviewer POST /pullRequest/reassign
func (h *PullRequestHandler) ReassignReviewer(c *gin.Context) {
	var req dto.ReassignPRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: err.Error(),
			},
		})
		return
	}

	response, err := h.prService.ReassignReviewer(c.Request.Context(), req)
	if err != nil {
		var serviceErr *services.ServiceError
		if errors.As(err, &serviceErr) {
			statusCode := http.StatusNotFound
			if serviceErr.Code == dto.ErrorCodePRMerged ||
				serviceErr.Code == dto.ErrorCodeNotAssigned ||
				serviceErr.Code == dto.ErrorCodeNoCandidate {
				statusCode = http.StatusConflict
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

	c.JSON(http.StatusOK, response)
}
