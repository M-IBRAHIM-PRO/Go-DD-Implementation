package booksvc

import (
	"context"
	"fmt"

	"library_system/internal/domain/book"
)

type Service struct {
	repo book.Repository
}

func NewService(repo book.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddBook(ctx context.Context, input AddBookInput) (*BookOutput, error) {
	isbn, err := book.NewISBN(input.ISBN)
	if err != nil {
		return nil, fmt.Errorf("invalid ISBN: %w", err)
	}
	b, err := book.NewBook(input.Title, input.Author, isbn)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, b); err != nil {
		return nil, fmt.Errorf("saving book: %w", err)
	}
	return toOutput(b), nil
}

func (s *Service) GetBook(ctx context.Context, id string) (*BookOutput, error) {
	b, err := s.repo.FindByID(ctx, book.BookID(id))
	if err != nil {
		return nil, err
	}
	return toOutput(b), nil
}

func toOutput(b *book.Book) *BookOutput {
	return &BookOutput{
		ID:      string(b.ID()),
		Title:   b.Title(),
		Author:  b.Author(),
		ISBN:    b.ISBN().String(),
		AddedAt: b.AddedAt().Format("2006-01-02"),
	}
}