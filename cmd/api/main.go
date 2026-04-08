package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	booksvc "library_system/internal/application/book"
	borrowingsvc "library_system/internal/application/borrowing"
	membersvc "library_system/internal/application/member"
	"library_system/internal/infrastructure/events"
	"library_system/internal/infrastructure/postgres"
	httphandler "library_system/internal/interfaces/http"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	log.Println("Connected to database")

	// Infrastructure
	bookRepo   := postgres.NewBookRepository(db)
	memberRepo := postgres.NewMemberRepository(db)
	borrowRepo := postgres.NewBorrowRepository(db)
	dispatcher := events.NewLogDispatcher()

	// Application services
	bookService   := booksvc.NewService(bookRepo)
	memberService := membersvc.NewService(memberRepo)
	borrowService := borrowingsvc.NewService(borrowRepo, bookRepo, memberRepo, dispatcher)

	// HTTP handlers
	bookHandler   := httphandler.NewBookHandler(bookService)
	memberHandler := httphandler.NewMemberHandler(memberService)
	borrowHandler := httphandler.NewBorrowHandler(borrowService)

	router := httphandler.NewRouter(bookHandler, memberHandler, borrowHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}