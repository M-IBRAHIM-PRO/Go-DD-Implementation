package borrowing

import "time"

// BookBorrowed — something that HAPPENED, immutable fact
type BookBorrowed struct {
    BorrowID  BorrowID
    BookID    string      // string, not book.BookID — no cross-domain imports
    MemberID  string      // string, not member.MemberID
    BorrowedAt time.Time
}

func (e BookBorrowed) EventName() string    { return "book.borrowed" }
func (e BookBorrowed) OccurredAt() time.Time { return e.BorrowedAt }

// BookReturned — another immutable fact
type BookReturned struct {
    BorrowID   BorrowID
    BookID     string
    MemberID   string
    ReturnedAt time.Time
}

func (e BookReturned) EventName() string    { return "book.returned" }
func (e BookReturned) OccurredAt() time.Time { return e.ReturnedAt }