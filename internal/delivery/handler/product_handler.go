package handler

import (
	"github.com/dettarune/kos-finder/internal/usecase"
	"github.com/sirupsen/logrus"
)

type KosHandler struct {
	KosUseCase *usecase.KosUseCase
	log     *logrus.Logger
}

func NewKosHandler(usecase *usecase.KosUseCase, log *logrus.Logger) *KosHandler {
	return &KosHandler{
		KosUseCase: usecase,
		log:     log,
	}
}