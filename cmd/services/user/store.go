package user

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/vickon16/go-jwt-mysql/cmd/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows, user)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == uuid.Nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) GetUserByID(id uuid.UUID) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows, user)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == uuid.Nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) CreateUser(user types.RegisterUserPayload) error {
	_, err := s.db.Exec("INSERT INTO users (id, firstName, lastName, email, password) VALUES (?, ?, ?, ?, ?)", user.ID, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoUser(rows *sql.Rows, user *types.User) (*types.User, error) {
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
