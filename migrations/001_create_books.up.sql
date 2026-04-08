CREATE TABLE books (
    id       VARCHAR(36)  PRIMARY KEY,
    title    VARCHAR(255) NOT NULL,
    author   VARCHAR(255) NOT NULL,
    isbn     VARCHAR(13)  NOT NULL UNIQUE,
    added_at TIMESTAMP    NOT NULL
);