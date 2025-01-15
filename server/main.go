package main

import (
	"context"
	"log"

	pb "faisal.com/bookProject/server/proto"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func main() {

	app := fiber.New()

	// Serve frontend files

	app.Static("/", "./client")

	// grpc connection
	// act as a grpc client :
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewBookServiceClient(conn)

	// handle form submission
	app.Post("/book", func(c *fiber.Ctx) error {

		// parse the book name from the client request
		var data struct {
			BookName string `json:"book_name"`
		}

		if err := c.BodyParser(&data); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		// call grpc AddBook rpc
		resp, err := client.AddBook(context.Background(), &pb.BookRequest{
			BookName: data.BookName,
		})

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to proccess book name")
		}

		return c.JSON(fiber.Map{
			"message": resp.Message,
		})

	})

	app.Get("/book", func(c *fiber.Ctx) error {
		// Call gRPC GetBooks method
		resp, err := client.GetBooks(context.Background(), &pb.EmptyRequest{})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch books")
		}

		return c.JSON(resp.Books)
	})

	log.Println("Fiber server is running on http://localhost:3000")
	app.Listen(":3000")

}
