package repository

import (
	"backend_golang/ent"
	"backend_golang/ent/transientmember"
	"backend_golang/internal/domain"
	"context"
	"fmt"
	"log"
)

type AuthRepository interface {
	CreateTransientMember(c context.Context, member *domain.TransientMember) (*domain.TransientMember, error)
	GetTransientMemberByID(c context.Context, id string) (*domain.TransientMember, error)
	CreateMember(c context.Context, member *domain.Member) (*domain.Member, error)
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
	var result *domain.TransientMember
	err := a.tx.WithTx(c, func(tx *ent.Tx) error {
		transientMember, err := tx.TransientMember.Create().
			SetTransientMemberID(member.ID).
			SetEmail(member.Email).
			SetPicture(member.Picture).
			SetNickname(member.Nickname).
			Save(c)
		if err != nil {
			return err
		}

		result = &domain.TransientMember{
			ID:       transientMember.TransientMemberID,
			Email:    transientMember.Email,
			Picture:  transientMember.Picture,
			Nickname: transientMember.Nickname,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *authRepository) GetTransientMemberByID(c context.Context, id string) (*domain.TransientMember, error) {
	var result *domain.TransientMember
	if err := a.tx.WithTx(c, func(tx *ent.Tx) error {
		member, err := tx.TransientMember.Query().Where(transientmember.TransientMemberID(id)).First(c)
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

func (a *authRepository) CreateMember(c context.Context, member *domain.Member) (*domain.Member, error) {

	transientMember, err := a.client.TransientMember.Query().Where(transientmember.TransientMemberID(member.ID)).First(c)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("transient member not found")
		}
		return nil, err
	}

	var result *domain.Member
	err = a.tx.WithTx(c, func(tx *ent.Tx) error {
		member, err := tx.Member.Create().
			SetMemberID(member.ID).
			SetEmail(transientMember.Email).
			SetPicture(transientMember.Picture).
			SetNickname(transientMember.Nickname).
			SetBio(member.Bio).
			SetPreferredRole(member.PreferredRole).
			Save(c)
		if err != nil {
			return err
		}

		result = &domain.Member{
			ID:            member.MemberID,
			Email:         transientMember.Email,
			Picture:       transientMember.Picture,
			Nickname:      transientMember.Nickname,
			Bio:           member.Bio,
			PreferredRole: member.PreferredRole,
		}

		return nil
	})
	return result, err
}
