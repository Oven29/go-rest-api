package services

import "go-rest-api/internal/api/dto"

type ServiceError struct {
	Code    dto.ErrorCode
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
