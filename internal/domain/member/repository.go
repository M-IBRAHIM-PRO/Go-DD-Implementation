package member

import "context"

type Repository interface {
    Save(ctx context.Context, member *Member) error
    FindByID(ctx context.Context, id MemberID) (*Member, error)
    FindByEmail(ctx context.Context, email Email) (*Member, error)
    ExistsByEmail(ctx context.Context, email Email) (bool, error)
}