package usecase

import (
	"context"
	"fmt"

	"github.com/dettarune/kos-finder/internal/exceptions"
	"github.com/dettarune/kos-finder/internal/model"
	"github.com/dettarune/kos-finder/internal/repository"
	"github.com/dettarune/kos-finder/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo      *repository.UserRepo
	validator *validator.Validate
	log       *logrus.Logger
	mail      *util.SMTPClient
	tokenUtil *util.TokenUtil
}

func NewUserUseCase(repo *repository.UserRepo, validator *validator.Validate, log *logrus.Logger, mail *util.SMTPClient, tokenUtil *util.TokenUtil) *UserUseCase {
	return &UserUseCase{
		repo:      repo,
		validator: validator,
		log:       log,
		mail:      mail,
		tokenUtil: tokenUtil,
	}
}

func (u *UserUseCase) Register(ctx context.Context, req *model.RegisterRequest) error {
	if err := u.validator.Struct(req); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			return exceptions.NewFailedValidationError(&verr)
		}
		u.log.Error("Validation parsing error:", err)
		return exceptions.NewBadRequestError("invalid request body")
	}

	existingUser, err := u.repo.FindUserByUsernameOrEmail(ctx, req.Username, req.Email)
	if err != nil {
		u.log.WithError(err).Error("failed to check existing user")
		return exceptions.NewInternalServerError()
	}

	if existingUser != nil {
		if err := util.CheckAuthConflict(existingUser, req); err != nil {
			return err
		}
	}

	hashedPW, err := util.HashPassword(req.Password)
	if err != nil {
		u.log.WithError(err).Error("failed to hash password")
		return exceptions.NewInternalServerError()
	}
	req.Password = string(hashedPW)

	if err := u.repo.InsertUser(ctx, req); err != nil {
		u.log.WithError(err).Error("failed to insert user")
		return exceptions.NewInternalServerError()
	}

	token, err := u.tokenUtil.CreateToken(&model.CreateToken{
		Username: req.Username,
	})
	if err != nil {
		u.log.WithError(err).Error("failed to create verification token")
		return exceptions.NewInternalServerError()
	}

	registerToken := fmt.Sprintf("http://localhost:2205/api/auth/verify?token=%s", token)

	if err := u.mail.SendMail("Verification Token", registerToken, req.Email); err != nil {
		u.log.WithError(err).Error("failed to send verification email")
		return exceptions.NewInternalServerError()
	}

	u.log.Infof("User %s registered successfully", req.Username)
	return nil
}

func (u *UserUseCase) VerifyEmail(ctx context.Context, token string) error {
    claims, err := u.tokenUtil.ParseToken(token)
    if err != nil {
        return err
    }

    username := claims.Username 

	u.log.Trace(username)

    err = u.repo.UpdateUserVerification(ctx, username, true)
    if err != nil {
        return err
    }

    return nil
}

func (u *UserUseCase) Login(ctx context.Context, reqUser *model.LoginRequest) (string, error) {
	if err := u.validator.Struct(reqUser); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			return "", exceptions.NewFailedValidationError(&verr)
		}
		u.log.Error("Validation parsing error:", err)
		return "", exceptions.NewBadRequestError("invalid request body")
	}

	user, err := u.repo.FindUserByUsernameOrEmail(ctx, reqUser.Username, reqUser.Username)
	if err != nil {
		u.log.WithError(err).Error("failed to find user")
		return "", exceptions.NewInternalServerError()
	}

	if user == nil {
		u.log.Warnf("User not found: %s", reqUser.Username)
		return "", exceptions.NewBadRequestError("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		u.log.Warnf("Invalid password for user: %s", reqUser.Username)
		return "", exceptions.NewBadRequestError("invalid credentials")
	}

	jwtToken, err := u.tokenUtil.CreateToken(&model.CreateToken{
		Username: user.Username,
	})
	if err != nil {
		u.log.WithError(err).Error("failed to create token")
		return "", exceptions.NewInternalServerError()
	}

	u.log.Infof("User logged in successfully: %s", reqUser.Username)
	return jwtToken, nil
}
