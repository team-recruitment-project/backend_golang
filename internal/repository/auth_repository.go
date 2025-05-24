package repository

import (
	"backend_golang/ent"
	"backend_golang/ent/transientmember"
	"backend_golang/internal/domain"
	"context"
	"log"
)

type AuthRepository interface {
	CreateTransientMember(c context.Context, member *domain.TransientMember) (*domain.TransientMember, error)
	GetTransientMemberByID(c context.Context, id string) (*domain.TransientMember, error)
}

type authRepository struct {
	client *ent.Client
	tx     *TransactionManager
}

func NewAuthRepository(client *ent.Client) AuthRepository {
	return &authRepository{
		client: client,
		tx:     NewTransactionManager(client),
	}
}

func (a *authRepository) CreateTransientMember(c context.Context, member *domain.TransientMember) (*domain.TransientMember, error) {
	panic("unimplemented")
}

func (a *authRepository) GetTransientMemberByID(c context.Context, id string) (*domain.TransientMember, error) {
	var result *domain.TransientMember
	if err := a.tx.WithTx(c, func(tx *ent.Tx) error {
		member, err := tx.TransientMember.Query().Where(transientmember.TransientMemberIDEQ(id)).Only(c)
		if err != nil {
			return err
		}

		result = &domain.TransientMember{
			ID:       member.TransientMemberID,
			Email:    member.Email,
			Picture:  member.Picture,
			Nickname: member.Nickname,
		}

		return nil
	}); err != nil {
		log.Printf("error getting transient member by id: %v", err)
		return nil, err
	}
	return result, nil
}
