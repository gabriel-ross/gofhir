package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

var PROJECT_ID string
var PORT string

func main() {
	var err error
	LoadConfigFromEnvironment()
	r := chi.NewRouter()
	ctx := context.Background()

	fsClient, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}
	defer fsClient.Close()

	clClient, err := gofhir.NewCloudLoggerClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}
	defer clClient.Close()

	// logger := clClient.Logger("default").StandardLogger(logging.Info)

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
