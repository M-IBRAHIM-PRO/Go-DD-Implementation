package borrowingsvc

type BorrowBookInput struct {
    BookID   string
    MemberID string
}

type ReturnBookInput struct {
    BorrowID string
}

type BorrowOutput struct {
    ID         string
    BookID     string
    MemberID   string
    Status     string
    BorrowedAt string
    DueAt      string
    ReturnedAt *string
    IsOverdue  bool
}