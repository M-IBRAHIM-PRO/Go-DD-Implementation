package membersvc

import (
    "context"
    "fmt"

    "library_system/internal/domain/member"
)

type Service struct {
    repo member.Repository // depends on the INTERFACE, not postgres
}

func NewService(repo member.Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) RegisterMember(ctx context.Context, input RegisterMemberInput) (*MemberOutput, error) {
    // 1. Build Value Objects (validation happens here)
    email, err := member.NewEmail(input.Email)
    if err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }

    // 2. Check business rule: no duplicate emails
    exists, err := s.repo.ExistsByEmail(ctx, email)
    if err != nil {
        return nil, fmt.Errorf("checking email existence: %w", err)
    }
    if exists {
        return nil, member.ErrAlreadyExists
    }

    // 3. Create the aggregate (domain enforces its own invariants)
    m, err := member.NewMember(input.Name, email)
    if err != nil {
        return nil, err
    }

    // 4. Persist
    if err := s.repo.Save(ctx, m); err != nil {
        return nil, fmt.Errorf("saving member: %w", err)
    }

    return toOutput(m), nil
}

func (s *Service) GetMember(ctx context.Context, id string) (*MemberOutput, error) {
    m, err := s.repo.FindByID(ctx, member.MemberID(id))
    if err != nil {
        return nil, err
    }
    return toOutput(m), nil
}

// Private helper — maps domain entity to output DTO
func toOutput(m *member.Member) *MemberOutput {
    return &MemberOutput{
        ID:       string(m.ID()),
        Name:     m.Name(),
        Email:    m.Email().String(),
        Status:   string(m.Status()),
        JoinedAt: m.JoinedAt().Format("2006-01-02"),
    }
}