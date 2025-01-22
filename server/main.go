package main

import (
	"context"
	"github.com/Trendyol/kafka-konsumer/v2"
	"google.golang.org/grpc/credentials/insecure"

	//"github.com/segmentio/kafka-go"
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
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	//// Create Kafka Writer for producing
	//kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
	//	Brokers: []string{"localhost:9092"}, // Kafka broker
	//	Topic:   "book-events",              // Kafka topic for updating books
	//})
	//defer kafkaWriter.Close()

	// Create kafka producer
	producer, _ := kafka.NewProducer(&kafka.ProducerConfig{
		Writer: kafka.WriterConfig{
			Brokers: []string{"localhost:9092"},
		},
	})

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

		//send a JSON-formatted response back to the client.
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

	app.Delete("/book", func(c *fiber.Ctx) error {
		var data struct {
			Id string `json:"id"`
		}
		if err := c.BodyParser(&data); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}
		resp, err := client.DeleteBook(context.Background(), &pb.BookDeletionRequest{
			Id: data.Id,
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete book")
		}

		return c.JSON(resp.Message)
	})

	app.Put("/book", func(c *fiber.Ctx) error {
		// Parse input
		var data struct {
			Id       string `json:"id"`
			BookName string `json:"book_name"`
		}

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		log.Printf("Received bookName to update: %s", data.BookName)

		// Produce update message to kafka
		//message := kafka.Message{
		//	Key:   []byte(data.Id),
		//	Value: []byte(data.BookName),
		//}

		const topicName = "book-events"
		_ = producer.Produce(context.Background(), kafka.Message{
			Topic: topicName,
			Key:   []byte(data.Id),
			Value: []byte(data.BookName),
		})

		if err != nil {
			log.Printf("Failed to send update message to Kafka: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to send update message to Kafka",
			})
		}

		log.Println("Message sent to Kafka")

		//if err := kafkaWriter.WriteMessages(context.Background(), message); err != nil {
		//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		//		"error": "Failed to send update message to Kafka",
		//	})
		//}

		return c.JSON(fiber.Map{
			"message": "Book updated successfully",
		})

	})

	log.Println("Fiber server is running on http://localhost:3000")
	app.Listen(":3000")

}
