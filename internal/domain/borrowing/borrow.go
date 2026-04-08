package borrowing

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "library_system/internal/domain"
)

type BorrowID string

func NewBorrowID() BorrowID {
    return BorrowID(uuid.New().String())
}

type BorrowStatus string

const (
    BorrowStatusActive    BorrowStatus = "active"
    BorrowStatusReturned  BorrowStatus = "returned"
)

// Borrow is the Aggregate Root of the Borrowing bounded context.
// KEY RULE: It references Book and Member only by their IDs (plain strings).
// This is how bounded contexts stay decoupled.
type Borrow struct {
    id         BorrowID
    bookID     string       // NOT book.BookID — just a string reference
    memberID   string       // NOT member.MemberID — just a string reference
    status     BorrowStatus
    borrowedAt time.Time
    returnedAt *time.Time   // pointer — nil until returned
    dueAt      time.Time

    // Events raised by this aggregate, collected and dispatched AFTER save
    events []domain.DomainEvent
}

// NewBorrow — creates a fresh borrowing record
// Notice: domain doesn't query DB here. The APPLICATION layer is responsible
// for checking if the book/member exist before calling this.
func NewBorrow(bookID, memberID string) (*Borrow, error) {
    if bookID == "" {
        return nil, errors.New("bookID cannot be empty")
    }
    if memberID == "" {
        return nil, errors.New("memberID cannot be empty")
    }

    now := time.Now()
    b := &Borrow{
        id:         NewBorrowID(),
        bookID:     bookID,
        memberID:   memberID,
        status:     BorrowStatusActive,
        borrowedAt: now,
        dueAt:      now.AddDate(0, 0, 14), // 2-week borrow period
    }

    // Raise the event — it gets dispatched after the aggregate is saved
    b.events = append(b.events, BookBorrowed{
        BorrowID:   b.id,
        BookID:     bookID,
        MemberID:   memberID,
        BorrowedAt: now,
    })

    return b, nil
}

func Reconstitute(id, bookID, memberID, status string, borrowedAt, dueAt time.Time, returnedAt *time.Time) *Borrow {
	return &Borrow{
		id:         BorrowID(id),
		bookID:     bookID,
		memberID:   memberID,
		status:     BorrowStatus(status),
		borrowedAt: borrowedAt,
		dueAt:      dueAt,
		returnedAt: returnedAt,
	}
}

// Return — behaviour method, enforces the state machine
func (b *Borrow) Return() error {
    if b.status == BorrowStatusReturned {
        return errors.New("book has already been returned")
    }
    now := time.Now()
    b.status = BorrowStatusReturned
    b.returnedAt = &now

    b.events = append(b.events, BookReturned{
        BorrowID:   b.id,
        BookID:     b.bookID,
        MemberID:   b.memberID,
        ReturnedAt: now,
    })

    return nil
}

func (b *Borrow) IsOverdue() bool {
    if b.status == BorrowStatusReturned {
        return false
    }
    return time.Now().After(b.dueAt)
}

// PopEvents — returns raised events and clears the internal list.
// Called by the application layer after a successful save.
func (b *Borrow) PopEvents() []domain.DomainEvent {
    evts := b.events
    b.events = nil
    return evts
}

// Getters
func (b *Borrow) ID() BorrowID          { return b.id }
func (b *Borrow) BookID() string        { return b.bookID }
func (b *Borrow) MemberID() string      { return b.memberID }
func (b *Borrow) Status() BorrowStatus  { return b.status }
func (b *Borrow) BorrowedAt() time.Time { return b.borrowedAt }
func (b *Borrow) ReturnedAt() *time.Time { return b.returnedAt }
func (b *Borrow) DueAt() time.Time      { return b.dueAt }