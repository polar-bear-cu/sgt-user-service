package models

type User struct {
	ID          string
	Email       string
	GoogleSub   string
	DisplayName string
	PictureURL  string
	Role        string
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)