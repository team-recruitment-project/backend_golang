package domain

type Team struct {
	ID          int
	Name        string
	Description string
	Headcount   int8
	CreatedBy   string
	Positions   []int
	Skills      []int
}
