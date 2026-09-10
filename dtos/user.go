package dtos

type UpdateProfileRequest struct {
	DisplayName string `json:"displayName" binding:"required"`
	PictureURL  string `json:"pictureUrl"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
}
