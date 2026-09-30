package models

import "time"

type User struct {
	ID             string
	Email          string
	GoogleSub      string
	DisplayName    string
	PictureURL     string
	Role           string
	CreatedAt      time.Time
	LastLoginAt    time.Time
	TimeInAdvanced int32
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
