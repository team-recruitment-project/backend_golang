package models

type Vacancy struct {
	Role    Role
	Vacancy int8
}

// NewVacancy creates a new Vacancy instance
func NewVacancy(role Role, vacancy int8) Vacancy {
	return Vacancy{
		Role:    role,
		Vacancy: vacancy,
	}
}
