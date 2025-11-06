package main

import (
	"net/http"

	"github.com/dettarune/kos-finder/db"
	"github.com/dettarune/kos-finder/internal/config"
	"github.com/dettarune/kos-finder/internal/delivery/handler"
	"github.com/dettarune/kos-finder/internal/middleware"
	"github.com/dettarune/kos-finder/internal/repository"
	"github.com/dettarune/kos-finder/internal/routes"
	"github.com/dettarune/kos-finder/internal/usecase"
	"github.com/dettarune/kos-finder/internal/util"
)

func main() {

	viper := config.NewViper()
	log := config.NewLogger(viper)
	validator := config.NewValidator(viper)

	db := db.NewDatabase(viper, log)

	tokenUtil := util.NewTokenUtils(viper) 

	smtpClient := util.NewSMTP(viper)

	userRepo := repository.NewUserRepo(db)
	productRepo := repository.NewProductRepo(db)
	
	userUseCase := usecase.NewUserUseCase(userRepo, validator, log, smtpClient, tokenUtil)
	kosUsecase := usecase.NewKosUseCase(productRepo, validator, log)

	UserHandler := handler.NewUserHandler(userUseCase, log)
	kosHandler := handler.NewKosHandler(kosUsecase, log)

	authmiddleware := middleware.NewAuthMiddleware(tokenUtil)

	router := routes.NewRouterConfig(UserHandler, kosHandler,authmiddleware)

	router.SetupGuestRoutes()
	router.SetupAuthRoutes()


	log.Info("db info : ", db)
	log.Info("\nServer started at http://localhost:2205")
	if err := http.ListenAndServe(":2205", router.Router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
