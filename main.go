package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-user-service/controllers"
	"github.com/polar-bear-cu/sgt-user-service/repositories"
	"github.com/polar-bear-cu/sgt-user-service/routes"
	"github.com/polar-bear-cu/sgt-user-service/usecases"
)

func main() {
	r := gin.Default()

	repo := repositories.NewInMemoryUser()
	uc := usecases.NewUser(repo)
	userCtrl := controllers.NewUser(uc)

	routes.Register(r, userCtrl)

	log.Println("listening :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
