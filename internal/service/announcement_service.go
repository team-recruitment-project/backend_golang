package service

import (
	"backend_golang/internal/service/models"
	"context"
)

type AnnouncementService interface {
	Announce(ctx context.Context, model *models.RegisterAnnouncement) error
}

type announcementService struct {
}

func NewAnnouncementController() AnnouncementService {
	return &announcementService{}
}

func (a *announcementService) Announce(ctx context.Context, model *models.RegisterAnnouncement) error {
	return nil
}
