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

const dummyUserID = "11111111-1111-1111-1111-111111111111"

func (ctl *UserController) GetMe(c *gin.Context) {
	user, err := ctl.uc.GetByID(c.Request.Context(), dummyUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

func (ctl *UserController) UpdateMe(c *gin.Context) {
	var req dtos.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := ctl.uc.UpdateProfile(c.Request.Context(), dummyUserID, req.DisplayName, req.PictureURL)
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
