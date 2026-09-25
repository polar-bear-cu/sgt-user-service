package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/polar-bear-cu/sgt-user-service/controllers"
	"github.com/polar-bear-cu/sgt-user-service/middlewares"
)

func Register(r *gin.Engine, user *controllers.UserController, jwtSecret string, swaggerEnabled bool) {
	r.GET("/health", controllers.GetHealth)

	if swaggerEnabled {
		r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	}

	v1 := r.Group("/api/v1")
	v1.Use(middlewares.RequireAuth(jwtSecret))

	v1.GET("/users/me", user.GetMe)
	v1.PATCH("/users/me", user.UpdateMe)
	v1.DELETE("/users/me", user.DeleteMe)

	//admin
	v1.GET("/users", user.GetAll)
	v1.GET("/users/:id", user.GetByID)
	v1.PATCH("/users/:id", user.UpdateRole)
	v1.DELETE("/users/:id", user.DeleteByID)
}
