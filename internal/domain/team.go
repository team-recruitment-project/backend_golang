package domain

type Team struct {
	ID          int
	Name        string
	Description string
	Headcount   int8
	Positions   []int
	Skills      []int
}
