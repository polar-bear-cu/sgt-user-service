package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-user-service/dtos"
	"github.com/polar-bear-cu/sgt-user-service/models"
	"github.com/polar-bear-cu/sgt-user-service/usecases"
)

type UserController struct {
	uc *usecases.UserUsecase
}

func NewUser(uc *usecases.UserUsecase) *UserController {
	return &UserController{uc: uc}
}

// GetMe godoc
// @Summary  get my profile
// @Tags     users
// @Produce  json
// @Success  200  {object}  dtos.UserResponse
// @Failure  404  {object}  map[string]string
// @Router   /api/v1/users/me [get]
func (ctl *UserController) GetMe(c *gin.Context) {
	user, err := ctl.uc.GetByID(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// UpdateMe godoc
// @Summary  update my profile
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    body  body      dtos.UpdateProfileRequest  true  "profile"
// @Success  200   {object}  dtos.UserResponse
// @Failure  400   {object}  map[string]string
// @Failure  500   {object}  map[string]string
// @Router   /api/v1/users/me [patch]
func (ctl *UserController) UpdateMe(c *gin.Context) {
	var req dtos.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := ctl.uc.UpdateProfile(c.Request.Context(), c.GetString("user_id"), req.DisplayName, req.PictureURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

func toUserResponse(u models.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		PictureURL:  u.PictureURL,
	}
}
