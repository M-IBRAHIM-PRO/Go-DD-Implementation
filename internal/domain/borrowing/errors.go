package borrowing

import "errors"

var (
    ErrNotFound         = errors.New("borrow record not found")
    ErrBookNotAvailable = errors.New("book is currently borrowed by someone else")
    ErrMemberSuspended  = errors.New("suspended members cannot borrow books")
)