package controllers

import (
	"errors"
	"net/http"
	"strconv"

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

// DeleteMe godoc
// @Summary  delete my own account
// @Tags     users
// @Success  204
// @Router   /api/v1/users/me [delete]
func (ctl *UserController) DeleteMe(c *gin.Context) {
	if err := ctl.uc.DeleteSelf(c.Request.Context(), c.GetString("user_id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetAll godoc
// @Summary  list all users
// @Tags     admin
// @Produce  json
// @Param    limit query int false "limit (default 10)"
// @Param    offset query int false "offset (default 0)"
// @Success  200 {array} dtos.UserResponse
// @Failure  403 {object} map[string]string
// @Router   /api/v1/users [get]
func (ctl *UserController) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := ctl.uc.GetAll(c.Request.Context(), c.GetString("user_id"), limit, offset)
	if err != nil {
		respondErr(c, err)
		return
	}

	resp := make([]dtos.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserResponse(u))
	}
	c.JSON(http.StatusOK, resp)
}

// GetByID godoc
// @Summary  get a user by id
// @Tags     admin
// @Produce  json
// @Param    id path string true "user id"
// @Success  200 {object} dtos.UserResponse
// @Failure  403 {object} map[string]string
// @Router   /api/v1/users/{id} [get]
func (ctl *UserController) GetByID(c *gin.Context) {
	user, err := ctl.uc.GetByIDAsAdmin(c.Request.Context(), c.GetString("user_id"), c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// UpdateRole godoc
// @Summary  change a user's role
// @Tags     admin
// @Accept   json
// @Produce  json
// @Param    id path string true "user id"
// @Param    body body dtos.UpdateRoleRequest true "new role"
// @Success  200 {object} dtos.UserResponse
// @Failure  403 {object} map[string]string
// @Router   /api/v1/users/{id} [patch]
func (ctl *UserController) UpdateRole(c *gin.Context) {
	var req dtos.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := ctl.uc.UpdateRole(c.Request.Context(), c.GetString("user_id"), c.Param("id"), req.Role)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// DeleteByID godoc
// @Summary  delete a user
// @Tags     admin
// @Param    id path string true "user id"
// @Success  204
// @Failure  403 {object} map[string]string
// @Router   /api/v1/users/{id} [delete]
func (ctl *UserController) DeleteByID(c *gin.Context) {
	if err := ctl.uc.DeleteUser(c.Request.Context(), c.GetString("user_id"), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func respondErr(c *gin.Context, err error) {
	if errors.Is(err, usecases.ErrForbidden) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func toUserResponse(u models.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		PictureURL:  u.PictureURL,
		Role:        u.Role,
	}
}
