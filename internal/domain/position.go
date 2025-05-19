package domain

import "backend_golang/internal/models"

type Position struct {
	// ID      int64
	Role    models.Role
	Vacancy int8
	// Team    int64
}
