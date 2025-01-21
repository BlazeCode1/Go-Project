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

func UpdateBook(id, newBookName string) error {
	_, err := Cluster.Bucket("books_bucket").DefaultCollection().Upsert(id, map[string]interface{}{
		"id":        id,
		"book_name": newBookName,
	}, nil)

	if err != nil {
		return err
	}

	return nil
}
