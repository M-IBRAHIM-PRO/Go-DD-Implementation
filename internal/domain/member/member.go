package member

import (
    "errors"
    "strings"
    "time"

    "github.com/google/uuid"
)

// MemberID — Value Object
type MemberID string

func NewMemberID() MemberID {
    return MemberID(uuid.New().String())
}

// Email — Value Object with its own validation
// This is a great example: Email isn't just a string,
// it has rules. The domain enforces those rules here.
type Email struct {
    value string
}

func NewEmail(raw string) (Email, error) {
    cleaned := strings.ToLower(strings.TrimSpace(raw))
    if !strings.Contains(cleaned, "@") || !strings.Contains(cleaned, ".") {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: cleaned}, nil
}

func (e Email) String() string { return e.value }

// MemberStatus — Value Object representing lifecycle state
type MemberStatus string

const (
    MemberStatusActive    MemberStatus = "active"
    MemberStatusSuspended MemberStatus = "suspended"
)

// Member — Aggregate Root
type Member struct {
    id        MemberID
    name      string
    email     Email
    status    MemberStatus
    joinedAt  time.Time
}

func NewMember(name string, email Email) (*Member, error) {
    if strings.TrimSpace(name) == "" {
        return nil, errors.New("name cannot be empty")
    }
    return &Member{
        id:       NewMemberID(),
        name:     name,
        email:    email,
        status:   MemberStatusActive, // new members always start active
        joinedAt: time.Now(),
    }, nil
}

func Reconstitute(id, name, email, status string, joinedAt time.Time) (*Member, error) {
	e, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	return &Member{
		id:       MemberID(id),
		name:     name,
		email:    e,
		status:   MemberStatus(status),
		joinedAt: joinedAt,
	}, nil
}

// Behaviour methods — domain logic lives ON the entity, not in a service
func (m *Member) Suspend() error {
    if m.status == MemberStatusSuspended {
        return errors.New("member is already suspended")
    }
    m.status = MemberStatusSuspended
    return nil
}

func (m *Member) Reactivate() error {
    if m.status == MemberStatusActive {
        return errors.New("member is already active")
    }
    m.status = MemberStatusActive
    return nil
}

func (m *Member) IsActive() bool {
    return m.status == MemberStatusActive
}

// Getters
func (m *Member) ID() MemberID       { return m.id }
func (m *Member) Name() string       { return m.name }
func (m *Member) Email() Email       { return m.email }
func (m *Member) Status() MemberStatus { return m.status }
func (m *Member) JoinedAt() time.Time { return m.joinedAt }