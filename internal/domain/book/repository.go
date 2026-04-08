package book

import "context"

// Repository is an INTERFACE defined in the domain layer.
// The domain says "I need something that can store and find books"
// but it doesn't care HOW (Postgres, MySQL, memory — doesn't matter).
type Repository interface {
    Save(ctx context.Context, book *Book) error
    FindByID(ctx context.Context, id BookID) (*Book, error)
}