package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ankush263/blog-api/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, arg *model.User) error {
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, password
	`
	return r.db.QueryRowContext(
		ctx, 
		query, 
		arg.Email, 
		arg.Password,
	).Scan(
		&arg.ID, 
		&arg.Email, 
		&arg.Password,
	)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error){
	query := `
		SELECT * FROM users
		WHERE email = $1
	`
	var u model.User

	err := r.db.QueryRowContext(
		ctx, 
		query, 
		email,
	).Scan(
		&u.ID, 
		&u.Email, 
		&u.Password,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("invalid credentials")
	}

	return &u, err
}
