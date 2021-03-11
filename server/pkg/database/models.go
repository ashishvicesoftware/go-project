package database

import (
	"time"
)


type User struct {
	ID                string
	Name              string
	Email             string
	IsActive          bool
	ImageUrl          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
