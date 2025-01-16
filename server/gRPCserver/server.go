package main

import (
	"context"
	"faisal.com/bookProject/couchbase"
	pb "faisal.com/bookProject/server/proto"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBookServiceServer
}

func (s *server) GetBooks(ctx context.Context, req *pb.EmptyRequest) (*pb.BookListResponse, error) {

	// Query all books details in Couchbase
	query := "SELECT id, book_name FROM `books_bucket`"
	rows, err := couchbase.Cluster.Query(query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %v", err)
	}

	var books []*pb.Book
	//  rows.Next() iterates through the result rows retrieved from a database query
	for rows.Next() {
		var row couchbase.Book
		if err := rows.Row(&row); err != nil {
			return nil, fmt.Errorf("failed to parse book row: %v", err)
		}
		books = append(books, &pb.Book{
			Id:       row.ID,
			BookName: row.BookName,
		})
	}

	return &pb.BookListResponse{
		Books: books,
	}, nil
}

// AddBook implements the gRPC method to handle book input
func (s *server) AddBook(ctx context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	log.Printf("Received bookName: %s", req.BookName)

	// Create a new Book instance
	book := couchbase.Book{
		ID:       uuid.New().String(), // Generate a new UUID for each book
		BookName: req.BookName,
	}

	// Insert the book into Couchbase
	if err := couchbase.InsertBook(book); err != nil {
		return nil, fmt.Errorf("failed to add book to Couchbase: %v", err)
	}

	// Return success message
	return &pb.BookResponse{
		Message: fmt.Sprintf("Book '%s' has been added successfully", req.BookName),
	}, nil
}

// DeleteBook implements the gRPC method to handle book deletion
func (s *server) DeleteBook(ctx context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	log.Printf("Received bookName to delete: %s", req.BookName)

	// Prepare a query to delete the book from Couchbase
	query := "DELETE FROM `books_bucket` WHERE book_name = $book_name"

	// Define query options with parameters
	options := &gocb.QueryOptions{
		NamedParameters: map[string]interface{}{
			// Replace $book_name with the value of req.BookName
			"book_name": req.BookName,
		},
	}

	// Execute the query
	_, err := couchbase.Cluster.Query(query, options)
	if err != nil {
		return nil, fmt.Errorf("failed to delete book from Couchbase: %v", err)
	}

	// Return success message
	return &pb.BookResponse{
		Message: fmt.Sprintf("Book '%s' has been deleted successfully", req.BookName),
	}, nil
}

func main() {

	// Initialize Couchbase
	couchbase.InitCouchbase("admin", "1q2w3e4r5t", "books_bucket")

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBookServiceServer(grpcServer, &server{})

	log.Println("Server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
