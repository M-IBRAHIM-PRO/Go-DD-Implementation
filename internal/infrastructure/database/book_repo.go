package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"library_system/internal/domain/book"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Save(ctx context.Context, b *book.Book) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO books (id, title, author, isbn, added_at) VALUES ($1, $2, $3, $4, $5)`,
		string(b.ID()), b.Title(), b.Author(), b.ISBN().String(), b.AddedAt(),
	)
	return err
}

func (r *BookRepository) FindByID(ctx context.Context, id book.BookID) (*book.Book, error) {
	var bid, title, author, isbn string
	var addedAt time.Time

	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, author, isbn, added_at FROM books WHERE id = $1`, string(id),
	).Scan(&bid, &title, &author, &isbn, &addedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, book.ErrNotFound
		}
		return nil, err
	}
	return book.Reconstitute(bid, title, author, isbn, addedAt)
}