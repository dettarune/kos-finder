package main

import (
	"net/http"

	"github.com/dettarune/kos-finder/internal/config"
	"github.com/dettarune/kos-finder/db"
	"github.com/dettarune/kos-finder/internal/delivery/handler"
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

	//util
	tokenUtil := util.NewTokenUtils(viper) 

	//smtp
	smtpClient := util.NewSMTP(viper)

	//repository
	userRepo := repository.NewUserRepo(db)
	productRepo := repository.NewProductRepo(db)
	
	//usecase
	userUseCase := usecase.NewUserUseCase(userRepo, validator, log, smtpClient, tokenUtil)
	productUseCase := usecase.NewProductUseCase(productRepo, validator, log)

	//handler
	UserHandler := handler.NewUserHandler(userUseCase, log)
	ProductHandler := handler.NewProductHandler(productUseCase, log)

	
	// router
	router := routes.NewRouterConfig(UserHandler, ProductHandler)

	router.SetupRoutes()

	log.Info("db info : ", db)
	log.Info("\nServer started at http://localhost:2205")
	if err := http.ListenAndServe(":2255", router.Router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
