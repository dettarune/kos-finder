package handler

import (
	"github.com/dettarune/kos-finder/internal/usecase"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	ProductUseCase *usecase.ProductUseCase
	log     *logrus.Logger
}

func NewProductHandler(usecase *usecase.ProductUseCase, log *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: usecase,
		log:     log,
	}
}