package borrowing

import "context"

type Repository interface {
    Save(ctx context.Context, borrow *Borrow) error
    Update(ctx context.Context, borrow *Borrow) error
    FindByID(ctx context.Context, id BorrowID) (*Borrow, error)
    FindActiveByMember(ctx context.Context, memberID string) ([]*Borrow, error)
    FindActiveByBook(ctx context.Context, bookID string) (*Borrow, error) // is a book currently borrowed?
}