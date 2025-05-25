package repository

type SkillRepository interface{}
type skillRepository struct{}

func NewSkillRepository() SkillRepository {
	return &skillRepository{}
}
