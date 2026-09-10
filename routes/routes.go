package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-user-service/controllers"
)

func Register(r *gin.Engine, user *controllers.UserController) {
	r.GET("/health", controllers.GetHealth)

	v1 := r.Group("/api/v1")
	v1.GET("/users/me", user.GetMe)
	v1.PATCH("/users/me", user.UpdateMe)
}
