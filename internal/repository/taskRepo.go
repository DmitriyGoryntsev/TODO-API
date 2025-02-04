package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/DmitriyGiryntsev/TODO-API/internal/models"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (t *TaskRepository) CreateNewTask(task *models.Task) error {
	stmt, err := t.DB.Prepare(`INSERT INTO tasks (userID, title, description, status, createdAt, updatedAt) VALUES ($1, $2, $3, $4, DEFAULT, NOW()) RETURNING id`)
	if err != nil {
		log.Print("cannot prepare statement to create new task:", err)
		return err
	}

	err = stmt.QueryRow(task.UserID, task.Title, task.Description, task.Status, task.Created_at, task.Updated_at).Scan(&task.ID)
	if err != nil {
		log.Print("cannot scan row to create new task:", err)
		return err
	}

	return nil
}

func (t *TaskRepository) GetAllTasksByUserID(userID int) ([]models.Task, error) {
	stmt, err := t.DB.Prepare(`SELECT id, userID, title, description, status, createdAt, updatedAt FROM tasks WHERE userID = $1`)
	if err != nil {
		log.Print("cannot prepare statement to get all tasks:", err)
		return nil, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		log.Print("cannot execute statement to get all tasks:", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Created_at, &task.Updated_at); err != nil {
			log.Print("cannot scan row to get all tasks:", err)
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *TaskRepository) GetTaskByID(taskID int, userID int) (*models.Task, error) {
	stmt, err := t.DB.Prepare(`SELECT id, userID, title, description, status, createdAt, updatedAt FROM tasks WHERE id = $1 AND userID = $2`)
	if err != nil {
		log.Print("cannot prepare statement to get task:", err)
		return nil, err
	}

	var task models.Task

	err = stmt.QueryRow(taskID, userID).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Created_at, &task.Updated_at)
	if err != nil {
		log.Print("cannot scan row to get task:", err)
		return nil, err
	}

	return &task, nil
}

func (t *TaskRepository) UpdateTask(task *models.Task) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4 AND userID = $5"

	result, err := t.DB.Exec(query, task.Title, task.Description, task.Status, task.ID, task.UserID)
	if err != nil {
		log.Print("cannot execute statement to update task:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("task not found or you don't have permission to update it")
	}

	return nil
}

func (t *TaskRepository) DeleteTask(taskID int, userID int) error {
	stmt, err := t.DB.Prepare("DELETE FROM tasks WHERE id = $1 AND userID = $2")
	if err != nil {
		log.Print("cannot prepare statement to delete task:", err)
		return err
	}

	result, err := stmt.Exec(taskID, userID)
	if err != nil {
		log.Print("cannot execute statement to delete task:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("task not found or you don't have permission to delete it")
	}

	return nil
}
