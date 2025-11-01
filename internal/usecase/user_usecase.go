package usecase

import (
	"context"
	"fmt"

	"github.com/dettarune/kos-finder/internal/entity"
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

func (u *UserUseCase) Register(ctx context.Context, req *entity.User) error {
	if err := u.validator.Struct(req); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}

	existingUser, err := u.repo.FindUserByUsernameOrEmail(ctx, req.Username, req.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		if err := util.CheckAuthConflict(existingUser, req); err != nil {
			return err
		}
	}

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(req.Password), 5)
	if err != nil {
		return err
	}
	req.Password = string(hashedPW)

	if err := u.repo.InsertUser(ctx, req); err != nil {
		return err
	}

	token, err := u.tokenUtil.CreateToken(req)
	if err != nil {
		return err
	}

	if err := u.mail.SendMail("Verification Token", token, req.Email); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}


func (u *UserUseCase) Login(ctx context.Context, reqUser *model.UserLogin) (string, error) {
	user, err := u.repo.FindUserByUsernameOrEmail(ctx, reqUser.Username, reqUser.Username)
	if err != nil {
		u.log.Error("Account Not Found:", reqUser.Username)
		return "", fmt.Errorf("account not found")
	}
	if user == nil {
		u.log.Error("User Not Found:", reqUser.Username)
		return "", fmt.Errorf("account not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		u.log.Error("Invalid Password for:", reqUser.Username)
		return "", fmt.Errorf("invalid password")
	}

	jwtToken, err := u.tokenUtil.CreateToken(user)
	if err != nil {
		u.log.Error("Failed to create JWT for:", reqUser.Username, "err:", err)
		return "", fmt.Errorf("failed to create token")
	}

	u.log.Info("User logged in:", reqUser.Username)
	return jwtToken, nil
}

// func (u *UserUseCase) VerifyAccount() (string, error) {

// }

func (u *UserUseCase) VerifyWithEmail(ctx context.Context, mail string) error{
	err := u.mail.SendMail("OIOIOOIOIO", "ini aja test anu, ", mail)
	if err != nil {
		return err
	}
	return nil
}

// func (u *UserUseCase) VerifyWithEmail(ctx context.Context, mail string) error{
// 	verifyToken := uuid.New()
// 	err := u.mail.SendMail("OIOIOOIOIO", "ini aja test anu, ", mail)
// 	if err != nil {
// 		return err
// 	}
// func (u *UserUseCase) VerifyWithEmail(ctx context.Context, mail string) error{
// 	verifyToken := uuid.New()
// 	err := u.mail.SendMail("OIOIOOIOIO", "ini aja test anu, ", mail)
// 	if err != nil {
// 		return err
// 	}
// func (u *UserUseCase) VerifyWithEmail(ctx context.Context, mail string) error{
// 	verifyToken := uuid.New()
// 	err := u.mail.SendMail("OIOIOOIOIO", "ini aja test anu, ", mail)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// 	return nil
// }

// 	return nil
// }
