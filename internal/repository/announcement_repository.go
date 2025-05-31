package repository

import (
	"backend_golang/ent"
	"backend_golang/ent/announcement"
	"backend_golang/ent/position"
	"backend_golang/ent/predicate"
	"backend_golang/ent/skill"
	"backend_golang/ent/team"
	"backend_golang/internal/domain"
	imodels "backend_golang/internal/models"
	"backend_golang/internal/service/models"
	"context"
	"log"
	"math"
)

type AnnouncementRepository interface {
	CreateAnnouncement(ctx context.Context, announcement models.RegisterAnnouncement) (*domain.Announcement, error)
	GetAnnouncement(ctx context.Context, announcementID int) (*domain.Announcement, error)
	GetLastAnnouncement(ctx context.Context, announcementID int) (*domain.Announcement, error)
	GetAnnouncements(ctx context.Context, page int, size int, skills []string, positions []string, keyword string) ([]domain.Announcement, error)
}

type announcementRepository struct {
	client *ent.Client
	tx     *TransactionManager
}

func NewAnnouncementRepository(client *ent.Client) AnnouncementRepository {
	return &announcementRepository{
		client: client,
		tx:     NewTransactionManager(client),
	}
}

func (a *announcementRepository) CreateAnnouncement(ctx context.Context, announcement models.RegisterAnnouncement) (*domain.Announcement, error) {
	var result *domain.Announcement
	err := a.tx.WithTx(ctx, func(tx *ent.Tx) error {
		announcement, err := tx.Announcement.Create().
			SetTitle(announcement.Title).
			SetContent(announcement.Content).
			SetTeamID(announcement.TeamID).
			Save(ctx)
		if err != nil {
			log.Printf("error creating announcement: %v", err)
			return err
		}
		result = &domain.Announcement{
			ID:        announcement.ID,
			Title:     announcement.Title,
			Content:   announcement.Content,
			CreatedAt: announcement.CreatedAt,
			UpdatedAt: announcement.UpdatedAt,
		}
		return nil
	})
	if err != nil {
		log.Printf("error getting announcement: %v", err)
		return nil, err
	}
	return result, nil
}

func (a *announcementRepository) GetAnnouncement(ctx context.Context, announcementID int) (*domain.Announcement, error) {
	announcement, err := a.client.Announcement.Query().
		WithTeam(
			func(tq *ent.TeamQuery) {
				tq.WithMembers().
					WithPositions().
					WithSkills()
			},
		).
		Where(announcement.ID(announcementID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	var skills []domain.Skill
	for _, skill := range announcement.Edges.Team.Edges.Skills {
		skills = append(skills, domain.Skill{
			Name: skill.Name,
		})
	}

	var positions []domain.Position
	for _, position := range announcement.Edges.Team.Edges.Positions {
		positions = append(positions, domain.Position{
			Role:    imodels.Role(position.Role),
			Vacancy: position.Vacancy,
		})
	}

	var members []domain.Member
	for _, member := range announcement.Edges.Team.Edges.Members {
		members = append(members, domain.Member{
			ID:            member.MemberID,
			Email:         member.Email,
			Nickname:      member.Nickname,
			Picture:       member.Picture,
			Bio:           member.Bio,
			PreferredRole: member.PreferredRole,
		})
	}

	return &domain.Announcement{
		ID:        announcement.ID,
		Title:     announcement.Title,
		Content:   announcement.Content,
		CreatedAt: announcement.CreatedAt,
		UpdatedAt: announcement.UpdatedAt,
		Team: &domain.Team{
			ID:          announcement.Edges.Team.ID,
			Name:        announcement.Edges.Team.Name,
			Description: announcement.Edges.Team.Description,
			Headcount:   announcement.Edges.Team.Headcount,
			CreatedBy:   announcement.Edges.Team.CreatedBy,
			Members:     members,
			Skills:      skills,
			Positions:   positions,
		},
	}, nil
}

func (a *announcementRepository) GetLastAnnouncement(ctx context.Context, announcementID int) (*domain.Announcement, error) {
	announcement, err := a.client.Announcement.Query().
		Order(ent.Desc(announcement.FieldCreatedAt)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Announcement{
		ID:        announcement.ID,
		Title:     announcement.Title,
		Content:   announcement.Content,
		CreatedAt: announcement.CreatedAt,
		UpdatedAt: announcement.UpdatedAt,
	}, nil
}

func (a *announcementRepository) GetAnnouncements(ctx context.Context, page int, size int, skills []string, positions []string, keyword string) ([]domain.Announcement, error) {
	query := a.client.Announcement.Query().WithTeam(
		func(tq *ent.TeamQuery) {
			tq.WithMembers().
				WithPositions().
				WithSkills()
		},
	)

	var conditions []predicate.Announcement
	if len(skills) > 0 {
		conditions = append(conditions, announcement.HasTeamWith(
			team.HasSkillsWith(skill.NameIn(skills...)),
		))
	}
	if len(positions) > 0 {
		conditions = append(conditions, announcement.HasTeamWith(
			team.HasPositionsWith(position.RoleIn(positions...)),
		))
	}

	if keyword != "" {
		conditions = append(conditions, announcement.Or(
			announcement.TitleContains(keyword),
			announcement.ContentContains(keyword),
		))
	}

	if len(conditions) > 0 {
		query = query.Where(announcement.And(conditions...))
	}

	announcements, err := query.
		Offset(int(math.Max(float64(page-1), 0)) * size).
		Limit(size).
		Order(ent.Desc(announcement.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var result []domain.Announcement
	for _, announcement := range announcements {
		var members []domain.Member
		for _, member := range announcement.Edges.Team.Edges.Members {
			members = append(members, domain.Member{
				ID:            member.MemberID,
				Email:         member.Email,
				Nickname:      member.Nickname,
				Picture:       member.Picture,
				Bio:           member.Bio,
				PreferredRole: member.PreferredRole,
			})
		}

		var skills []domain.Skill
		for _, skill := range announcement.Edges.Team.Edges.Skills {
			skills = append(skills, domain.Skill{
				Name: skill.Name,
			})
		}

		var positions []domain.Position
		for _, position := range announcement.Edges.Team.Edges.Positions {
			positions = append(positions, domain.Position{
				Role:    imodels.Role(position.Role),
				Vacancy: position.Vacancy,
			})
		}

		result = append(result, domain.Announcement{
			ID:        announcement.ID,
			Title:     announcement.Title,
			Content:   announcement.Content,
			CreatedAt: announcement.CreatedAt,
			UpdatedAt: announcement.UpdatedAt,
			Team: &domain.Team{
				ID:          announcement.Edges.Team.ID,
				Name:        announcement.Edges.Team.Name,
				Description: announcement.Edges.Team.Description,
				Headcount:   announcement.Edges.Team.Headcount,
				CreatedBy:   announcement.Edges.Team.CreatedBy,
				Members:     members,
				Skills:      skills,
				Positions:   positions,
			},
		})
	}
	return result, nil
}
