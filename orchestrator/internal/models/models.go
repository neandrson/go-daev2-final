package models

import (
	"time"

	"github.com/neandrson/go-daev2-final/orchestrator/internal/services/calculation"
)

type Expression struct {
	ID         int               `json:"id"`
	Status     string            `json:"status"`
	Result     float64           `json:"result"`
	UserID     int               `json:"-"`
	BinaryTree *calculation.Tree `json:"-"`
}

type Task struct {
	ID            int           `json:"id"`
	ExpressionID  int           `json:"-"`
	Status        string        `json:"-"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
}

type User struct {
	ID           int
	Login        string
	PasswordHash string
	Password     string
}
