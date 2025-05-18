package models

type Vacancy struct {
	Role    Role
	Vacancy int
}

// NewVacancy creates a new Vacancy instance
func NewVacancy(role Role, vacancy int) Vacancy {
	return Vacancy{
		Role:    role,
		Vacancy: vacancy,
	}
}
