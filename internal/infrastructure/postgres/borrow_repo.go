package postgres

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "library_system/internal/domain/borrowing"
)

type BorrowRepository struct {
    db *sql.DB
}

func NewBorrowRepository(db *sql.DB) *BorrowRepository {
    return &BorrowRepository{db: db}
}

func (r *BorrowRepository) Save(ctx context.Context, b *borrowing.Borrow) error {
    _, err := r.db.ExecContext(ctx, `
        INSERT INTO borrowings (id, book_id, member_id, status, borrowed_at, due_at, returned_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
        string(b.ID()), b.BookID(), b.MemberID(),
        string(b.Status()), b.BorrowedAt(), b.DueAt(), b.ReturnedAt(),
    )
    return err
}

func (r *BorrowRepository) Update(ctx context.Context, b *borrowing.Borrow) error {
    _, err := r.db.ExecContext(ctx, `
        UPDATE borrowings SET status = $1, returned_at = $2 WHERE id = $3
    `, string(b.Status()), b.ReturnedAt(), string(b.ID()))
    return err
}

func (r *BorrowRepository) FindByID(ctx context.Context, id borrowing.BorrowID) (*borrowing.Borrow, error) {
    row := r.db.QueryRowContext(ctx,
        `SELECT id, book_id, member_id, status, borrowed_at, due_at, returned_at
         FROM borrowings WHERE id = $1`, string(id))
    return scanBorrow(row)
}

func (r *BorrowRepository) FindActiveByBook(ctx context.Context, bookID string) (*borrowing.Borrow, error) {
    row := r.db.QueryRowContext(ctx,
        `SELECT id, book_id, member_id, status, borrowed_at, due_at, returned_at
         FROM borrowings WHERE book_id = $1 AND status = 'active'`, bookID)
    return scanBorrow(row)
}

func (r *BorrowRepository) FindActiveByMember(ctx context.Context, memberID string) ([]*borrowing.Borrow, error) {
    rows, err := r.db.QueryContext(ctx,
        `SELECT id, book_id, member_id, status, borrowed_at, due_at, returned_at
         FROM borrowings WHERE member_id = $1 AND status = 'active'`, memberID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []*borrowing.Borrow
    for rows.Next() {
        var id, bookID, memberID, status string
		var borrowedAt, dueAt time.Time
		var returnedAt *time.Time
		if err := rows.Scan(&id, &bookID, &memberID, &status, &borrowedAt, &dueAt, &returnedAt); err != nil {
			return nil, err
		}
		result = append(result, borrowing.Reconstitute(id, bookID, memberID, status, borrowedAt, dueAt, returnedAt))
	}
    return result, nil
}

func scanBorrow(row *sql.Row) (*borrowing.Borrow, error) {
    var id, bookID, memberID, status string
    var borrowedAt, dueAt time.Time
    var returnedAt *time.Time

    if err := row.Scan(&id, &bookID, &memberID, &status, &borrowedAt, &dueAt, &returnedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, borrowing.ErrNotFound
        }
        return nil, err
    }
    return borrowing.Reconstitute(id, bookID, memberID, status, borrowedAt, dueAt, returnedAt), nil
}
