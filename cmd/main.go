package main

import (
	"context"
	"fmt"

	"github.com/niiilov/go-dog-trapping/internal/application"
	"github.com/niiilov/go-dog-trapping/internal/config"
	"github.com/niiilov/go-dog-trapping/internal/dto"
	"github.com/niiilov/go-dog-trapping/internal/repository"
	"github.com/niiilov/go-dog-trapping/internal/service"
	"github.com/niiilov/go-dog-trapping/pkg/jwt"
	"github.com/niiilov/go-dog-trapping/pkg/postgres"
)

// @title Go Dog Trapping API
// @version 1.0
// @description API сервиса для отлова бродячих собак.

// @host localhost:8091
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	ctx := context.Background()

	privateKey, err := jwt.LoadPrivateKey("./certs/private.pem")
	if err != nil {

		//логирование
	}

	publicKey, err := jwt.LoadPublicKey("./certs/public.pem")
	if err != nil {

		//логирование
	}

	config := config.NewConfig()

	jwtService := jwt.NewServiceJWT(privateKey, publicKey, dto.RefreshTimeExpr, dto.AccesTimeExpr)

	pg, err := postgres.NewPostgres(ctx, config.Postgres)
	if err != nil {

		//логирование
	}
	fmt.Println(pg.Ping(ctx))

	repo := repository.New(pg)

	service := service.New(repo)

	handlers := application.NewHandlers(service, jwtService)

	application.StartApplication(":8091", handlers)
}
