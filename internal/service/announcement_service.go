package service

import (
	"backend_golang/internal/repository"
	"backend_golang/internal/service/models"
	"context"
)

type AnnouncementService interface {
	Announce(ctx context.Context, model models.RegisterAnnouncement) (int, error)
}

type announcementService struct {
	announcementRepository repository.AnnouncementRepository
}

func NewAnnouncementController(announcementRepository repository.AnnouncementRepository) AnnouncementService {
	return &announcementService{
		announcementRepository: announcementRepository,
	}
}

func (a *announcementService) Announce(ctx context.Context, model models.RegisterAnnouncement) (int, error) {
	announcement, err := a.announcementRepository.CreateAnnouncement(ctx, model)
	if err != nil {
		return 0, err
	}
	return announcement.ID, nil
}
