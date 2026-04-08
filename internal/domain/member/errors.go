package member

import "errors"

// Sentinel errors — callers can check: errors.Is(err, member.ErrNotFound)
// This keeps error handling explicit and domain-specific
var (
    ErrNotFound      = errors.New("member not found")
    ErrAlreadyExists = errors.New("member with this email already exists")
)