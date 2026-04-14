package store

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err

		}
	}
	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"` // - means we are going to ignore the value in the struct
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

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(string) (*User, error)
	UpdateUser(*User) error
}

func (pgus *PostgresUserStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password_hash, bio) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := pgus.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (pgus *PostgresUserStore) GetUserByUsername(username string) (*User, error) {

	user := &User{
		PasswordHash: password{},
	}

	query := `SELECT id, username, email, password_hash, bio, created_at, updated_at FROM users WHERE username = $1`

	err := pgus.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, nil
}

func (pgus *PostgresUserStore) UpdateUser(user *User) error {
	query := `UPDATE users
	SET username = $1, email = $2, bio = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4
	RETURNING updated_at`

	err := pgus.db.QueryRow(query, user.Username, user.Email, user.Bio, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
