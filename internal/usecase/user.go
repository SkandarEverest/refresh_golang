package usecase

import (
	"context"

	db "github.com/SkandarEverest/refresh-golang/db/sqlc"
	"github.com/SkandarEverest/refresh-golang/internal/dto"
	"github.com/SkandarEverest/refresh-golang/internal/exception"

	"github.com/sirupsen/logrus"
)

type UserUseCase struct {
	DB  *db.SQLStore
	Log *logrus.Logger
}

func NewUserUseCase(db *db.SQLStore, logger *logrus.Logger) *UserUseCase {
	return &UserUseCase{
		DB:  db,
		Log: logger,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *dto.UserRequest) (*db.User, error) {
	arg := db.CreateUserParams{
		Username:       "admin",
		HashedPassword: "admin",
		Fullname:       "administrator",
		Email:          request.Email,
	}

	user, err := c.DB.CreateUser(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, exception.Conflict("unique violation")
		}
		return nil, exception.ServerError(err.Error())
	}

	return &user, nil
}
