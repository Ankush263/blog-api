package repository

import (
	"context"
	"database/sql"

	"github.com/Ankush263/blog-api/internal/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, arg *model.Post) error {
	query := `
		INSERT INTO posts (title, content)
		VALUES ($1, $2)
		RETURNING id, title, content
	`
	return r.db.QueryRowContext(
		ctx, 
		query, 
		arg.Title, 
		arg.Content,
	).Scan(
		&arg.ID, 
		&arg.Title, 
		&arg.Content,
	)
}

func (r *PostRepository) GetAll(ctx context.Context) ([]model.Post, error) {
	query := `
		SELECT * FROM posts
		ORDER BY updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.ID,
			&p.CreatedAt,
			&p.Title,
			&p.Content,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetById(ctx context.Context, id int) (*model.Post, error) {
	query := `
		SELECT * FROM posts p
		WHERE p.id = ?
	`

	var p model.Post
	err := r.db.QueryRowContext(
		ctx, 
		query, 
		id,
	).Scan(
		&p.ID,
		&p.CreatedAt,
		&p.Title,
		&p.Content,
	)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepository) Update(ctx context.Context, arg *model.Post, id int) error {
	query := `
		UPDATE posts
		SET title = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, arg.Title, arg.Content, id)

	return err
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,`
		DELETE FROM posts WHERE id = $1
	`, id)

	return err
}
