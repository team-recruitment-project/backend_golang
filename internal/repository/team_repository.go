package repository

import (
	"backend_golang/ent"
	"backend_golang/ent/member"
	"backend_golang/ent/position"
	"backend_golang/ent/skill"
	"backend_golang/ent/team"
	"backend_golang/internal/domain"
	"backend_golang/internal/models"
	"context"
	"log"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, createTeam *domain.Team) (*domain.Team, error)
	DeleteTeam(ctx context.Context, teamID int) error
	FindByID(ctx context.Context, teamID int) (*domain.Team, error)
	JoinTeam(ctx context.Context, teamID int, memberID string) error
}

type teamRepository struct {
	client *ent.Client
	tx     *TransactionManager
}

func NewTeamRepository(client *ent.Client) TeamRepository {
	return &teamRepository{
		client: client,
		tx:     NewTransactionManager(client),
	}
}

func (t *teamRepository) CreateTeam(ctx context.Context, createTeam *domain.Team) (*domain.Team, error) {
	// TODO : createTeam 을 서비스의 dto 가 아니라 리포지토리단의 domain 모델로 변경
	var result *domain.Team
	err := t.tx.WithTx(ctx, func(tx *ent.Tx) error {

		// 技術スタックがあるか探し、なければ作成する
		var skills []*ent.Skill
		for _, s := range createTeam.Skills {
			// TODO skill repository に移行するのが良いかも
			foundSkill, err := tx.Skill.Query().Where(skill.Name(s.Name)).First(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					// TODO skill repository に移行するのが良いかも
					foundSkill, err = tx.Skill.Create().
						SetName(s.Name).
						Save(ctx)
					if err != nil {
						log.Printf("error creating skill: %v", err)
						return err
					}
				} else {
					return err
				}
			}
			skills = append(skills, foundSkill)
		}

		positions := []*ent.Position{}
		for _, vacancy := range createTeam.Positions {
			savedPosition, err := tx.Position.Create().
				SetRole(string(vacancy.Role)).
				SetVacancy(vacancy.Vacancy).
				Save(ctx)
			if err != nil {
				return err
			}
			positions = append(positions, savedPosition)
		}

		foundMember, err := t.client.Member.Query().Where(member.MemberID(createTeam.CreatedBy)).First(ctx)
		if err != nil {
			log.Printf("error finding member: %v", err)
			return err
		}

		team, err := tx.Team.Create().
			SetName(createTeam.Name).
			SetDescription(createTeam.Description).
			SetHeadcount(createTeam.Headcount).
			SetCreatedBy(createTeam.CreatedBy).
			AddMembers(foundMember).
			AddPositions(positions...).
			AddSkills(skills...).
			Save(ctx)
		if err != nil {
			return err
		}

		// メンバーとチームを紐づける
		member, err := tx.Member.Query().Where(member.MemberID(createTeam.CreatedBy)).First(ctx)
		if err != nil {
			return err
		}
		_, err = member.Update().SetTeams(team).Save(ctx)
		if err != nil {
			log.Printf("error updating member: %v", err)
			return err
		}

		relatedPositions := make([]domain.Position, len(team.Edges.Positions))
		for i, position := range team.Edges.Positions {
			relatedPositions[i] = domain.Position{
				Role:    models.Role(position.Role),
				Vacancy: position.Vacancy,
			}
		}

		relatedSkills := make([]domain.Skill, len(team.Edges.Skills))
		for i, skill := range team.Edges.Skills {
			relatedSkills[i] = domain.Skill{
				Name: skill.Name,
			}
		}

		result = &domain.Team{
			ID:          team.ID,
			Name:        team.Name,
			Description: team.Description,
			Headcount:   team.Headcount,
			CreatedBy:   team.CreatedBy,
			Positions:   relatedPositions,
			Skills:      relatedSkills,
		}
		return nil
	})
	if err != nil {
		log.Printf("error creating team: %v", err)
		return nil, err
	}
	return result, nil
}

func (t *teamRepository) DeleteTeam(ctx context.Context, teamID int) error {
	return t.tx.WithTx(ctx, func(tx *ent.Tx) error {
		// First, delete all positions associated with the team
		_, err := tx.Position.Delete().Where(
			position.HasTeamWith(team.ID(teamID)),
		).Exec(ctx)
		if err != nil {
			return err
		}

		// Then delete the team
		err = tx.Team.DeleteOneID(teamID).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (t *teamRepository) FindByID(ctx context.Context, teamID int) (*domain.Team, error) {
	team, err := t.client.Team.Query().Where(team.ID(teamID)).WithMembers().WithPositions().WithSkills().First(ctx)
	if err != nil {
		return nil, err
	}

	members := make([]domain.Member, len(team.Edges.Members))
	for i, member := range team.Edges.Members {
		members[i] = domain.Member{
			ID:            member.MemberID,
			Email:         member.Email,
			Picture:       member.Picture,
			Nickname:      member.Nickname,
			Bio:           member.Bio,
			PreferredRole: member.PreferredRole,
		}
	}

	positions := make([]domain.Position, len(team.Edges.Positions))
	for i, position := range team.Edges.Positions {
		positions[i] = domain.Position{
			Role:    models.Role(position.Role),
			Vacancy: position.Vacancy,
		}
	}
	skills := make([]domain.Skill, len(team.Edges.Skills))
	for i, skill := range team.Edges.Skills {
		skills[i] = domain.Skill{
			Name: skill.Name,
		}
	}

	return &domain.Team{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		Headcount:   team.Headcount,
		CreatedBy:   team.CreatedBy,
		Members:     members,
		Positions:   positions,
		Skills:      skills,
	}, nil
}

func (t *teamRepository) JoinTeam(ctx context.Context, teamID int, memberID string) error {
	return t.tx.WithTx(ctx, func(tx *ent.Tx) error {
		// 멤버 조회
		memberEnt, err := tx.Member.Query().Where(member.MemberID(memberID)).First(ctx)
		if err != nil {
			log.Printf("error finding member: %v", err)
			return err
		}

		// 팀 조회
		teamEnt, err := tx.Team.Query().Where(team.ID(teamID)).First(ctx)
		if err != nil {
			log.Printf("error finding team: %v", err)
			return err
		}

		// 포지션 조회 (role, teamID로)
		posEnt, err := tx.Position.Query().
			Where(
				position.RoleEQ(memberEnt.PreferredRole),
				position.HasTeamWith(team.ID(teamID)),
			).
			ForUpdate(). // lock 획득
			First(ctx)
		if err != nil {
			log.Printf("error finding position: %v", err)
			return err
		}

		// vacancy 감소
		_, err = posEnt.Update().
			SetVacancy(posEnt.Vacancy - 1).
			Save(ctx)
		if err != nil {
			log.Printf("error updating position: %v", err)
			return err
		}

		// 멤버와 팀 연결 (1 : 1 관계)
		_, err = memberEnt.Update().
			SetTeamsID(teamEnt.ID).
			Save(ctx)
		if err != nil {
			log.Printf("error updating member: %v", err)
			return err
		}

		return nil
	})
}
