package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang-stady/internal/models"
	"log"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

const tableName = `public.tasks`

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db}
}

func (ts *TaskStore) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `SELECT * FROM ` + tableName + ` ORDER BY created_at DESC`

	rows, err := ts.db.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("error getting tasks: %w", err)
	}

	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		v := models.Task{}
		err := rows.StructScan(&v)
		if err != nil {
			return nil, fmt.Errorf("error scanning task: %w", err)
		}

		tasks = append(tasks, v)
	}

	return tasks, nil
}

func (ts *TaskStore) GetById(id int) (*models.Task, error) {
	var tasks models.Task

	query := strings.Builder{}
	query.WriteString(`SELECT * FROM ` + tableName + ` WHERE id = $1`)

	err := ts.db.Get(&tasks, query.String(), id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("task with %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting task: %w", err)
	}

	return &tasks, nil
}

func (ts *TaskStore) Create(task models.CreateTaskInput) (*models.Task, error) {
	var v models.Task

	query := strings.Builder{}
	query.WriteString(`
INSERT INTO ` + tableName + ` (payload) VALUES ($1)
RETURNING id, payload, completed, created_at, updated_at
`)

	payload, err := json.Marshal(task.Payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	err = ts.db.QueryRowx(query.String(), payload).StructScan(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (ts *TaskStore) UpdateCompleted(id int, task models.UpdateTaskInput) (*models.Task, error) {
	var v models.Task

	query := strings.Builder{}
	query.WriteString(`
UPDATE ` + tableName + ` SET (completed = $1) WHERE id = $2
RETURNING id, payload, completed, created_at, updated_at
`)

	err := ts.db.Get(&v, query.String(), task.Completed, id)
	if err != nil {
		return nil, fmt.Errorf("error updating task: %w", err)
	}

	if reflect.DeepEqual(v, models.Task{}) {
		return nil, fmt.Errorf("task with id %d not found", id)
	}

	return &v, nil
}

func (ts *TaskStore) Delete(id int) error {
	query := strings.Builder{}
	query.WriteString(`DELETE FROM ` + tableName + ` WHERE id = $1`)

	_, err := ts.db.Exec(query.String(), id)

	return err
}
