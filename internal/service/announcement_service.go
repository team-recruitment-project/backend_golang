package service

import (
	"backend_golang/ent"
	"backend_golang/internal/models"
	"backend_golang/internal/repository"
	imodels "backend_golang/internal/service/models"
	"context"
	"errors"
	"time"
)

type AnnouncementService interface {
	Announce(ctx context.Context, model imodels.RegisterAnnouncement) (int, error)
	GetAnnouncement(ctx context.Context, announcementID int) (*imodels.AnnouncementResponse, error)
}

type announcementService struct {
	announcementRepository repository.AnnouncementRepository
	teamRepository         repository.TeamRepository
}

func NewAnnouncementService(announcementRepository repository.AnnouncementRepository, teamRepository repository.TeamRepository) AnnouncementService {
	return &announcementService{
		announcementRepository: announcementRepository,
		teamRepository:         teamRepository,
	}
}

func (a *announcementService) Announce(ctx context.Context, model imodels.RegisterAnnouncement) (int, error) {
	// チーム長であることを確認
	team, err := a.teamRepository.FindByID(ctx, model.TeamID)
	if err != nil {
		return 0, err
	}
	if team.CreatedBy != model.MemberID {
		return 0, errors.New("you are not the team leader")
	}

	// 最新のアナウンスを取得
	lastAnnouncement, err := a.announcementRepository.GetLastAnnouncement(ctx, model.TeamID)
	if err != nil && !ent.IsNotFound(err) {
		return 0, err
	}

	// アナウンスが存在し、24時間経過していない場合はエラー
	if lastAnnouncement != nil && !lastAnnouncement.CreatedAt.Add(24*time.Hour).Before(time.Now()) {
		return 0, errors.New("you can only announce once every 24 hours")
	}

	// アナウンスを作成
	announcement, err := a.announcementRepository.CreateAnnouncement(ctx, model)
	if err != nil {
		return 0, err
	}
	return announcement.ID, nil
}

func (a *announcementService) GetAnnouncement(ctx context.Context, announcementID int) (*imodels.AnnouncementResponse, error) {
	announcement, err := a.announcementRepository.GetAnnouncement(ctx, announcementID)

	var vacancies []models.Vacancy
	for _, position := range announcement.Team.Positions {
		vacancies = append(vacancies, models.Vacancy{
			Role:    position.Role,
			Vacancy: position.Vacancy,
		})
	}
	if err != nil {
		return nil, err
	}

	return &imodels.AnnouncementResponse{
		ID:        announcement.ID,
		Title:     announcement.Title,
		Content:   announcement.Content,
		CreatedAt: announcement.CreatedAt,
		UpdatedAt: announcement.UpdatedAt,
		Team: &imodels.TeamResponse{
			ID:          announcement.Team.ID,
			Name:        announcement.Team.Name,
			Description: announcement.Team.Description,
			Headcount:   announcement.Team.Headcount,
			CreatedBy:   announcement.Team.CreatedBy,
			Members:     announcement.Team.Members,
			Vacancies:   vacancies,
			Skills:      announcement.Team.Skills,
		},
	}, nil
}
