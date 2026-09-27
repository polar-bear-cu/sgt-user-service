package dtos

type UpdateProfileRequest struct {
	DisplayName string `json:"displayName" binding:"required"`
	PictureURL  string `json:"pictureUrl"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	Role        string `json:"role"`
}
