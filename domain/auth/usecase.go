package auth

import (
	"github.com/helyus1412/auth-service/dto"
	"github.com/helyus1412/auth-service/model"
	"github.com/helyus1412/auth-service/pkg/utils"
)

type Usecase interface {
	Register(*dto.RegisterRequest) error
}

type usecase struct {
	repository Repository
}

func NewUsecase(repository Repository) Usecase {
	return &usecase{repository}
}

func (u *usecase) Register(payload *dto.RegisterRequest) error {
	// hash password

	hashedPassword, err := utils.HashPassword(payload.Password, 12)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    payload.Email,
		Password: hashedPassword,
	}

	err = u.repository.Insert(user)

	if err != nil {
		return err
	}

	return nil
}
