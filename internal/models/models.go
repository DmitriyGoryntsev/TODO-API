package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username" validate:"required,min=3,max=50"`
	Email      string    `json:"email" validate:"required,email"`
	Password   string    `json:"password" validate:"required,min=8"`
	Role       string    `json:"role" validate:"oneof=admin user"`
	Created_at time.Time `json:"created_at"`
}

type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"required,min=10,max=500"`
	Status      string    `json:"status" validate:"oneof=pending in_progress completed"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}
