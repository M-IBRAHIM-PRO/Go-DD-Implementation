CREATE TABLE members (
    id        VARCHAR(36)  PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    email     VARCHAR(255) NOT NULL UNIQUE,
    status    VARCHAR(20)  NOT NULL DEFAULT 'active',
    joined_at TIMESTAMP    NOT NULL
);