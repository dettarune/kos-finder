package usecase

import (
	"github.com/dettarune/kos-finder/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type KosUseCase struct {
	UserRepo      *repository.ProductRepo
	validator *validator.Validate
	log       *logrus.Logger
}

func NewKosUseCase(repo *repository.ProductRepo, validator *validator.Validate, log *logrus.Logger, ) *KosUseCase {
	return &KosUseCase{
		UserRepo: repo,
		validator: validator,
		log: log,
	}
}

// func (s *KosUseCase) CreateKos() (*model.CreateKosResponse, error) {

// 	err := A.Validate.Struct(req)


// 	return nil, nil
// }