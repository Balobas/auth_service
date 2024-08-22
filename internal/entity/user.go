package entity

import "time"

type User struct {
	Id              int64
	Name            string
	Email           string
	Role            string
	Password        string
	ConfirmPassword string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
