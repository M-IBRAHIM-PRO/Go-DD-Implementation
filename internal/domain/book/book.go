package book

import (
    "errors"
    "strings"
    "time"

    "github.com/google/uuid"
)

// BookID is a Value Object — it's just a typed string, but it has meaning
type BookID string

func NewBookID() BookID {
    return BookID(uuid.New().String())
}

// ISBN is a Value Object with its own validation rules
type ISBN struct {
    value string
}

func NewISBN(raw string) (ISBN, error) {
    cleaned := strings.ReplaceAll(raw, "-", "")
    if len(cleaned) != 10 && len(cleaned) != 13 {
        return ISBN{}, errors.New("ISBN must be 10 or 13 digits")
    }
    return ISBN{value: cleaned}, nil
}

func (i ISBN) String() string { return i.value }

// Book is our Aggregate Root (the main entity)
type Book struct {
    id        BookID
    title     string
    author    string
    isbn      ISBN
    addedAt   time.Time
}

// Constructor enforces invariants — a Book can NEVER be in an invalid state
func NewBook(title, author string, isbn ISBN) (*Book, error) {
    if strings.TrimSpace(title) == "" {
        return nil, errors.New("title cannot be empty")
    }
    if strings.TrimSpace(author) == "" {
        return nil, errors.New("author cannot be empty")
    }
    return &Book{
        id:      NewBookID(),
        title:   title,
        author:  author,
        isbn:    isbn,
        addedAt: time.Now(),
    }, nil
}

func Reconstitute(id, title, author, isbn string, addedAt time.Time) (*Book, error) {
	i, err := NewISBN(isbn)
	if err != nil {
		return nil, err
	}
	return &Book{
		id:      BookID(id),
		title:   title,
		author:  author,
		isbn:    i,
		addedAt: addedAt,
	}, nil
}

// Getters — fields are private, domain controls access
func (b *Book) ID() BookID        { return b.id }
func (b *Book) Title() string     { return b.title }
func (b *Book) Author() string    { return b.author }
func (b *Book) ISBN() ISBN        { return b.isbn }
func (b *Book) AddedAt() time.Time { return b.addedAt }