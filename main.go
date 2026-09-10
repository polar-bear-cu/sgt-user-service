package main

import (
	"context"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userv1 "github.com/polar-bear-cu/sgt-proto/gen/go/user/v1"
	"github.com/polar-bear-cu/sgt-user-service/config"
	"github.com/polar-bear-cu/sgt-user-service/controllers"
	_ "github.com/polar-bear-cu/sgt-user-service/docs"
	grpcserver "github.com/polar-bear-cu/sgt-user-service/grpc"
	"github.com/polar-bear-cu/sgt-user-service/repositories"
	"github.com/polar-bear-cu/sgt-user-service/routes"
	"github.com/polar-bear-cu/sgt-user-service/usecases"
)

// @title        User Service API
// @version      1.0
// @description  REST API for the Subglutee user service
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

	repo := repositories.NewUserPostgres(pool)
	uc := usecases.NewUser(repo)
	userCtrl := controllers.NewUser(uc)

	go startGRPC(ctx, cfg.GRPCPort, uc)

	r := gin.Default()
	routes.Register(r, userCtrl, cfg.SwaggerEnabled)

	log.Println("listening :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func startGRPC(ctx context.Context, port string, uc *usecases.UserUsecase) {
	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	gs := grpclib.NewServer()
	userv1.RegisterUserServiceServer(gs, grpcserver.NewUserServer(uc))
	reflection.Register(gs)

	log.Println("gRPC :" + port)
	if err := gs.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
