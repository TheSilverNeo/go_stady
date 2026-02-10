package models

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID        int             `json:"id" db:"id"`
	Payload   json.RawMessage `json:"payload" db:"payload"`
	Completed bool            `json:"completed" db:"completed"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type CreateTaskInput struct {
	Payload   json.RawMessage `json:"payload"`
	Completed bool            `json:"completed"`
}

type UpdateTaskInput struct {
	Id        *int  `json:"id"`
	Completed *bool `json:"completed"`
}
