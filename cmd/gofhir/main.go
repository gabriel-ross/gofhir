package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

var PROJECT_ID string
var PORT string

func main() {
	var err error
	LoadConfigFromEnvironment()
	r := chi.NewRouter()
	ctx := context.Background()

	// fsClient, err := firestore.NewClient(ctx, PROJECT_ID)
	// if err != nil {
	// 	log.Fatalf("Failed to create client: %v", err)
	// 	return
	// }
	// defer fsClient.Close()

	// clClient, err := gofhir.NewCloudLoggerClient(ctx, PROJECT_ID)
	// if err != nil {
	// 	log.Fatalf("Failed to create client: %v", err)
	// 	return
	// }
	// defer clClient.Close()

	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "requests",
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}

	err = kafkaWriter.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte("some key"),
			Value: []byte("some value"),
		},
	)
	if err != nil {
		log.Fatalf("failed to write messages: %v", err)
		return
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err = kafkaWriter.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(""),
				Value: []byte(r.RemoteAddr),
			},
		)
		if err != nil {
			log.Fatalf("failed to write request to kafka: %v", err)
		}
	})

	http.ListenAndServe(":"+PORT, r)
}

func LoadConfigFromEnvironment() {
	godotenv.Load(".env")
	PROJECT_ID = os.Getenv("PROJECT_ID")
	PORT = os.Getenv("PORT")

	// Default value if not set
	if PORT == "" {
		PORT = "8080"
	}
}

type Config struct {
	PROJECT_ID string `env:"PROJECT_ID" required:"true" default:"-"`
	PORT       string `env:"PORT" required:"false" default:"8080"`
}
