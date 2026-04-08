package borrowingsvc

import (
    "context"
    "fmt"

    "library_system/internal/domain"
    "library_system/internal/domain/book"
    "library_system/internal/domain/borrowing"
    "library_system/internal/domain/member"
)

type Service struct {
    borrowRepo borrowing.Repository
    bookRepo   book.Repository    // used for existence check only
    memberRepo member.Repository  // used for existence + status check only
    dispatcher domain.EventDispatcher
}

func NewService(
    borrowRepo borrowing.Repository,
    bookRepo book.Repository,
    memberRepo member.Repository,
    dispatcher domain.EventDispatcher,
) *Service {
    return &Service{
        borrowRepo: borrowRepo,
        bookRepo:   bookRepo,
        memberRepo: memberRepo,
        dispatcher: dispatcher,
    }
}

func (s *Service) BorrowBook(ctx context.Context, input BorrowBookInput) (*BorrowOutput, error) {
    // 1. Validate member exists and is active
    m, err := s.memberRepo.FindByID(ctx, member.MemberID(input.MemberID))
    if err != nil {
        return nil, fmt.Errorf("member lookup: %w", err)
    }
    if !m.IsActive() {
        return nil, borrowing.ErrMemberSuspended
    }

    // 2. Validate book exists
    _, err = s.bookRepo.FindByID(ctx, book.BookID(input.BookID))
    if err != nil {
        return nil, fmt.Errorf("book lookup: %w", err)
    }

    // 3. Check book is not already borrowed (domain rule)
    existing, err := s.borrowRepo.FindActiveByBook(ctx, input.BookID)
    if err != nil && err != borrowing.ErrNotFound {
        return nil, fmt.Errorf("checking book availability: %w", err)
    }
    if existing != nil {
        return nil, borrowing.ErrBookNotAvailable
    }

    // 4. Create the aggregate — domain raises the BookBorrowed event internally
    borrow, err := borrowing.NewBorrow(input.BookID, input.MemberID)
    if err != nil {
        return nil, err
    }

    // 5. Persist
    if err := s.borrowRepo.Save(ctx, borrow); err != nil {
        return nil, fmt.Errorf("saving borrow: %w", err)
    }

    // 6. Dispatch events AFTER successful save
    s.dispatcher.Dispatch(borrow.PopEvents())

    return toOutput(borrow), nil
}

func (s *Service) ReturnBook(ctx context.Context, input ReturnBookInput) (*BorrowOutput, error) {
    // 1. Find the borrow record
    borrow, err := s.borrowRepo.FindByID(ctx, borrowing.BorrowID(input.BorrowID))
    if err != nil {
        return nil, err
    }

    // 2. Apply domain behaviour — Return() enforces state machine
    if err := borrow.Return(); err != nil {
        return nil, err
    }

    // 3. Persist the updated aggregate
    if err := s.borrowRepo.Update(ctx, borrow); err != nil {
        return nil, fmt.Errorf("updating borrow: %w", err)
    }

    // 4. Dispatch events
    s.dispatcher.Dispatch(borrow.PopEvents())

    return toOutput(borrow), nil
}

func toOutput(b *borrowing.Borrow) *BorrowOutput {
    out := &BorrowOutput{
        ID:         string(b.ID()),
        BookID:     b.BookID(),
        MemberID:   b.MemberID(),
        Status:     string(b.Status()),
        BorrowedAt: b.BorrowedAt().Format("2006-01-02"),
        DueAt:      b.DueAt().Format("2006-01-02"),
        IsOverdue:  b.IsOverdue(),
    }
    if b.ReturnedAt() != nil {
        s := b.ReturnedAt().Format("2006-01-02")
        out.ReturnedAt = &s
    }
    return out
}