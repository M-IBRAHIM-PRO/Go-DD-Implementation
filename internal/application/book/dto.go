package booksvc

type AddBookInput struct {
	Title  string
	Author string
	ISBN   string
}

type BookOutput struct {
	ID      string
	Title   string
	Author  string
	ISBN    string
	AddedAt string
}