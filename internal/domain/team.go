package domain

type Team struct {
	ID          int
	Name        string
	Description string
	Headcount   int8
	CreatedBy   string
	Members     []Member
	Positions   []Position
	Skills      []Skill
}
