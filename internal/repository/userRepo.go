package repository

import (
	"database/sql"
	"log"

	"github.com/DmitriyGiryntsev/TODO-API/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (u *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User

	stmt, err := u.DB.Prepare("SELECT id, username, email, password, role, created_at FROM users WHERE id = $1")
	if err != nil {
		log.Print("cannot prepare statement to get user:", err)
		return nil, err
	}

	err = stmt.QueryRow(id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.Created_at)
	if err != nil {
		log.Print("cannot scan row to get user:", err)
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) CreateNewUser(user *models.User) error {
	stmt, err := u.DB.Prepare("INSERT INTO users (username, email, password, role, createdAt) VALUES ($1, $2, $3, $4, DEFAULT) RETURNING id, createdAt")
	if err != nil {
		log.Print("cannot prepare statement to create new user:", err)
		return err
	}

	_, err = stmt.Exec(user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		log.Print("cannot execute statement to create new user:", err)
		return err
	}

	return nil
}

func (u *UserRepository) UpdateUser(user *models.User) error {
	stmt, err := u.DB.Prepare("UPDATE users SET username = $1, email = $2, password = $3, role = $4, createdAt = $5 WHERE id = $6")
	if err != nil {
		log.Print("cannot prepare statement to update user:", err)
		return err
	}

	_, err = stmt.Exec(user.Username, user.Email, user.Password, user.Role, user.Created_at, user.ID)
	if err != nil {
		log.Print("cannot execute statement to update user:", err)
		return err
	}

	return nil
}

func (u *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	stmt, err := u.DB.Prepare("SELECT id, username, email, password, role, createdAt FROM users WHERE email = $1")
	if err != nil {
		log.Print("cannot prepare statement to get user:", err)
		return nil, err
	}

	err = stmt.QueryRow(email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.Created_at)
	if err != nil {
		log.Print("cannot scan row to get user:", err)
		return nil, err
	}

	return &user, nil
}
