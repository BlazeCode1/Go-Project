package couchbase

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

type Book struct {
	ID       string `json:"id"`
	BookName string `json:"book_name"`
}

func InsertBook(book Book) error {
	collection := GetCollection()
	_, err := collection.Upsert(book.ID, book, &gocb.UpsertOptions{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("could not insert book: %v", err)
	}
	return nil
}
