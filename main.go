package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-user-service/config"
	"github.com/polar-bear-cu/sgt-user-service/controllers"
	"github.com/polar-bear-cu/sgt-user-service/repositories"
	"github.com/polar-bear-cu/sgt-user-service/routes"
	"github.com/polar-bear-cu/sgt-user-service/usecases"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := config.ConnectPostgres(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	r := gin.Default()

	repo := repositories.NewUserPostgres(pool)
	uc := usecases.NewUser(repo)
	userCtrl := controllers.NewUser(uc)

	routes.Register(r, userCtrl)

	log.Println("listening :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
