package store

import (
	"database/sql"
	"log"
	"time"
)

type password struct {
	plainText *string
	hash []byte
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password    `json:"-"` // - means we are going to ignore the value in the struct
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db     *sql.DB
	logger *log.Logger
}

func NewPostgresUserStore(db *sql.DB, logger *log.Logger) *PostgresUserStore {
	return &PostgresUserStore{db: db, logger: logger}
}

type UserStore interface{
	CreateUser(*User) error
	GetUserByUsername(string) (*User, error)
	UpdateUser(*User) error
}

