package usecase

import (
	"context"
	"errors"

	db "github.com/SkandarEverest/refresh-golang/db/sqlc"
	"github.com/SkandarEverest/refresh-golang/pkg/bcrypt"
	jwt "github.com/SkandarEverest/refresh-golang/pkg/jwt"

	"github.com/SkandarEverest/refresh-golang/internal/dto"
	"github.com/SkandarEverest/refresh-golang/internal/exception"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

type UserUseCase struct {
	DB     *db.SQLStore
	Log    *logrus.Logger
	Config *viper.Viper
}

func NewUserUseCase(db *db.SQLStore, logger *logrus.Logger, config *viper.Viper) *UserUseCase {
	return &UserUseCase{
		DB:     db,
		Log:    logger,
		Config: config,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *dto.AuthRequest) (*string, error) {
	hashedPassword := bcrypt.HashPassword(request.Password)
	arg := db.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
	}

	user, err := c.DB.CreateUser(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, exception.Conflict("unique violation")
		}
		return nil, exception.ServerError(err.Error())
	}

	token := jwt.CreateToken(user.ID, user.Email, []byte(c.Config.GetString("JWT_SECRET")))

	return &token, nil
}

func (c *UserUseCase) Login(ctx context.Context, request *dto.AuthRequest) (*string, error) {

	user, err := c.DB.GetUserFromEmail(ctx, request.Email)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, exception.NotFound("User not found")
		}
		return nil, exception.ServerError(err.Error())
	}

	err = bcrypt.ComparePassword(request.Password, user.HashedPassword)
	if err != nil {
		return nil, exception.BadRequest("password is wrong")
	}

	token := jwt.CreateToken(user.ID, user.Email, []byte(c.Config.GetString("JWT_SECRET")))

	return &token, nil
}

func (c *UserUseCase) GetUser(ctx context.Context, userid int64) (*db.User, error) {

	user, err := c.DB.GetUser(ctx, userid)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, exception.NotFound("User not found")
		}
		return nil, exception.ServerError(err.Error())
	}

	return &user, nil
}

func (c *UserUseCase) UpdateUser(ctx context.Context, request *dto.UpdateUserRequest, userid int64) (*db.User, error) {

	arg := db.UpdateUserParams{
		ID:              userid,
		Email:           request.Email,
		Username:        request.Username,
		UserImageUri:    request.UserImageUri,
		CompanyName:     request.CompanyName,
		CompanyImageUri: request.CompanyImageUri,
	}

	user, err := c.DB.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, exception.NotFound("User not found")
		}
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, exception.Conflict("unique violation")
		}
		return nil, exception.ServerError(err.Error())
	}

	return &user, nil
}
