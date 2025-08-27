package main

import (
	"github.com/niiilov/go-dog-trapping/internal/application"
	"github.com/niiilov/go-dog-trapping/internal/dto"
	"github.com/niiilov/go-dog-trapping/internal/service"
	"github.com/niiilov/go-dog-trapping/pkg/jwt"
)

func main() {

	privateKey, err := jwt.LoadPrivateKey("./certs/private.pem")
	if err != nil {

		//логирование
	}

	publicKey, err := jwt.LoadPublicKey("./certs/public.pem")
	if err != nil {

		//логирование
	}

	jwtService := jwt.NewServiceJWT(privateKey, publicKey, dto.RefreshTimeExpr, dto.AccesTimeExpr)

	service := service.New()

	handlers := application.NewHandlers(service, jwtService)

	application.StartApplication(":8091", handlers)
}
