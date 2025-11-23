package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-rest-api/internal/api/dto"
	"go-rest-api/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// SetIsActive POST /users/setIsActive
func (h *UserHandler) SetIsActive(c *gin.Context) {
	var req dto.SetIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: err.Error(),
			},
		})
		return
	}

	user, err := h.userService.SetIsActive(c.Request.Context(), req)
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

	c.JSON(http.StatusOK, dto.SetIsActiveResponse{
		User: *user,
	})
}

// GetUserReviews GET /users/getReview?user_id=...
func (h *UserHandler) GetUserReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    dto.ErrorCodeNotFound,
				Message: "user_id query parameter is required",
			},
		})
		return
	}

	response, err := h.userService.GetUserReviews(c.Request.Context(), userID)
	if err != nil {
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
