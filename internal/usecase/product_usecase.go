package usecase

import (
	"github.com/dettarune/kos-finder/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ProductUseCase struct {
	repo      *repository.ProductRepo
	validator *validator.Validate
	log       *logrus.Logger
}

func NewProductUseCase(repo *repository.ProductRepo, validator *validator.Validate, log *logrus.Logger, ) *ProductUseCase {
	return &ProductUseCase{
		repo: repo,
		validator: validator,
		log: log,
	}
}