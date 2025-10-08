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
	Login(context.Context, *dto.LoginRequest) utils.Result
	ListUser(context.Context) utils.Result
	Edit(context.Context, *dto.EditRequest) utils.Result
	Delete(context.Context, int64) utils.Result
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

func (u *usecase) Login(ctx context.Context, payload *dto.LoginRequest) (result utils.Result) {
	// hash password
	user, err := u.repository.GetByEmail(payload.Email)
	if err != nil {
		result.Error = httpError.NewInternalServerError(err.Error())
		return result
	}

	if user == nil {
		result.Error = httpError.NewBadRequest("user not found")
		return result
	}

	// check password
	if err := utils.CheckPassword(user.Password, payload.Password); err != nil {
		result.Error = httpError.NewBadRequest("invalid password")
		return result
	}

	return result
}

func (u *usecase) ListUser(ctx context.Context) (result utils.Result) {
	users, err := u.repository.GetAll()
	if err != nil {
		result.Error = httpError.NewInternalServerError(err.Error())
		return result
	}

	result.Data = users

	return result
}

func (u *usecase) Edit(ctx context.Context, payload *dto.EditRequest) (result utils.Result) {
	// check if user exists
	user, err := u.repository.GetByID(payload.ID)
	if err != nil {
		result.Error = httpError.NewInternalServerError(err.Error())
		return result
	}

	if user == nil {
		result.Error = httpError.NewBadRequest("user not found")
		return result
	}

	userPayload := &model.User{
		ID:    payload.ID,
		Email: payload.Email,
	}
	// hash password
	if payload.Password != "" {
		hashedPassword, err := utils.HashPassword(payload.Password, 12)
		if err != nil {
			result.Error = httpError.NewInternalServerError(err.Error())
			return result
		}

		userPayload.Password = hashedPassword
	}

	if payload.Email == "" {
		userPayload.Email = user.Email
	}

	if payload.Password == "" {
		userPayload.Password = user.Password
	}
	err = u.repository.Update(userPayload)

	if err != nil {
		result.Error = httpError.NewInternalServerError(err.Error())
		return result
	}

	return result
}

func (u *usecase) Delete(ctx context.Context, id int64) (result utils.Result) {
	// check if user exists
	user, err := u.repository.GetByID(id)
	if err != nil {
		result.Error = httpError.NewInternalServerError(err.Error())
		return result
	}

	if user == nil {
		result.Error = httpError.NewBadRequest("user not found")
		return result
	}

	err = u.repository.SoftDelete(id)
	if err != nil {
		result.Error = httpError.NewInternalServerError(err.Error())
		return result
	}

	return result
}
