package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"library_system/internal/domain/member"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) Save(ctx context.Context, m *member.Member) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO members (id, name, email, status, joined_at) VALUES ($1, $2, $3, $4, $5)`,
		string(m.ID()), m.Name(), m.Email().String(), string(m.Status()), m.JoinedAt(),
	)
	return err
}

func (r *MemberRepository) FindByID(ctx context.Context, id member.MemberID) (*member.Member, error) {
	return r.scanOne(r.db.QueryRowContext(ctx,
		`SELECT id, name, email, status, joined_at FROM members WHERE id = $1`, string(id),
	))
}

func (r *MemberRepository) FindByEmail(ctx context.Context, email member.Email) (*member.Member, error) {
	return r.scanOne(r.db.QueryRowContext(ctx,
		`SELECT id, name, email, status, joined_at FROM members WHERE email = $1`, email.String(),
	))
}

func (r *MemberRepository) ExistsByEmail(ctx context.Context, email member.Email) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM members WHERE email = $1)`, email.String(),
	).Scan(&exists)
	return exists, err
}

func (r *MemberRepository) scanOne(row *sql.Row) (*member.Member, error) {
	var id, name, email, status string
	var joinedAt time.Time
	if err := row.Scan(&id, &name, &email, &status, &joinedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, member.ErrNotFound
		}
		return nil, err
	}
	return member.Reconstitute(id, name, email, status, joinedAt)
}