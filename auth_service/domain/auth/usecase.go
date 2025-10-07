package auth

import (
	"context"

	"github.com/helyus1412/auth-service/dto"
	"github.com/helyus1412/auth-service/model"
	httpError "github.com/helyus1412/auth-service/pkg/httpError"
	"github.com/helyus1412/auth-service/pkg/utils"
)

type Usecase interface {
	Register(context.Context, *dto.RegisterRequest) utils.Result
}

type usecase struct {
	repository Repository
}

func NewUsecase(repository Repository) Usecase {
	return &usecase{repository}
}

func (u *usecase) Register(ctx context.Context, payload *dto.RegisterRequest) (result utils.Result) {
	// hash password

	hashedPassword, err := utils.HashPassword(payload.Password, 12)
	if err != nil {
		result.Error = httpError.NewInternalServerError(err.Error())
		return result
	}

	user := &model.User{
		Email:    payload.Email,
		Password: hashedPassword,
	}

	err = u.repository.Insert(user)

	if err != nil {
		result.Error = httpError.NewInternalServerError(err.Error())
		return result
	}

	return result
}
