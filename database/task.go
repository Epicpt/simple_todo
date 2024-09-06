package database

import (
	"database/sql"
	"errors"

	"main.go/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

// AddTask добавляет задачу в базу данных, возвращает id задачи
func (r Repository) AddTask(task model.Task) (int64, error) {

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTasks получает последние 50 задач из базы данных
func (r Repository) GetTasks() ([]model.Task, error) {
	var tasks []model.Task

	rows, err := r.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask получает задачу по id
func (r Repository) GetTask(id string) (model.Task, error) {
	var task model.Task

	if err := r.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return task, err
		}

	}
	return task, nil
}

// UpdateTask обновляет все поля задачи
func (r Repository) UpdateTask(task model.Task) error {

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("expected to affect 1 row")
	}
	return nil
}

// UpdateTask обновляет дату задачи
func (r Repository) UpdateTaskDate(task model.Task) error {

	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := r.db.Exec(query, task.Date, task.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return err
	}
	return nil
}

// DeleteTask удаляет задачу из базы данных
func (r Repository) DeleteTask(id string) error {

	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	res, err := r.db.Exec(deleteQuery, id)
	if err != nil {

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("expected to affect 1 row")
	}

	return nil
}
