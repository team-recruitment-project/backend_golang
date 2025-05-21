package models

type Vacancy struct {
	Role    Role `json:"role"`
	Vacancy int8 `json:"vacancy"`
}

// NewVacancy creates a new Vacancy instance
func NewVacancy(role Role, vacancy int8) Vacancy {
	return Vacancy{
		Role:    role,
		Vacancy: vacancy,
	}
}
