package repository

import (
	"backend_golang/ent"
	"backend_golang/internal/domain"
	"backend_golang/internal/service/models"
	"context"
)

type AnnouncementRepository interface {
	CreateAnnouncement(ctx context.Context, announcement models.RegisterAnnouncement) (*domain.Announcement, error)
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
			Save(ctx)
		if err != nil {
			return err
		}
		result = &domain.Announcement{
			ID:      announcement.ID,
			Title:   announcement.Title,
			Content: announcement.Content,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
