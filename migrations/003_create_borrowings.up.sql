CREATE TABLE borrowings (
    id          VARCHAR(36) PRIMARY KEY,
    book_id     VARCHAR(36) NOT NULL REFERENCES books(id),
    member_id   VARCHAR(36) NOT NULL REFERENCES members(id),
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    borrowed_at TIMESTAMP   NOT NULL,
    due_at      TIMESTAMP   NOT NULL,
    returned_at TIMESTAMP   NULL
);

CREATE INDEX idx_borrowings_book_active   ON borrowings(book_id)   WHERE status = 'active';
CREATE INDEX idx_borrowings_member_active ON borrowings(member_id) WHERE status = 'active';