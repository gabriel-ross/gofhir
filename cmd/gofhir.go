package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir"
	"github.com/gabriel-ross/gofhir/interceptor"
	"github.com/gabriel-ross/gofhir/patient"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

var APPLICATION_URL string
var PROJECT_ID string
var PORT string
var KAFKA_BROKER_ADDRESS string

func main() {
	var err error
	LoadConfigFromEnvironment()
	ctx := context.Background()

	// Instantiate dependencies
	r := chi.NewRouter()
	fsClient, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Println("Failed to create firestore client: ", err)
	}
	kafkaClient := gofhir.NewKafkaWriter(KAFKA_BROKER_ADDRESS)
	logger, err := gofhir.NewCloudLogger(ctx, PROJECT_ID, "final-demo")
	if err != nil {
		log.Println("Failed to create logging client: ", err)
	}

	defer func() {
		fsClient.Close()
		kafkaClient.Close()
		logger.Close()
	}()

	// Instantiate services and register interceptors
	p := patient.New(r, fsClient, APPLICATION_URL+":"+PORT, "patients")
	p.RegisterInterceptors(
		&interceptor.StreamingInterceptor{
			KafkaWriter: *kafkaClient,
		},
		&interceptor.BenchmarkInterceptor{
			RequestTime:     map[string]time.Time{},
			PreDbQueryTime:  map[string]time.Time{},
			PostDbQueryTime: map[string]time.Time{},
			ResponseTime:    map[string]time.Time{},
		},
		&interceptor.LoggingInterceptor{
			Logger: logger,
		},
	)

	log.Println("Starting server on port: " + PORT)
	logger.Log("Starting server on port: " + PORT)
	http.ListenAndServe(":"+PORT, r)
	log.Println("shutting down server...")
}

func LoadConfigFromEnvironment() {
	godotenv.Load(".env")
	APPLICATION_URL = os.Getenv("APPLICATION_URL")
	PROJECT_ID = os.Getenv("PROJECT_ID")
	PORT = os.Getenv("PORT")
	KAFKA_BROKER_ADDRESS = os.Getenv("KAFKA_BROKER_ADDRESS")

	// Default value if not set
	if PORT == "" {
		PORT = "8080"
	}
}

type Config struct {
	EnvFilePath     string
	APPLICATION_URL string `env:"APPLICATION_URL" required:"true" default:"localhost"`
	PROJECT_ID      string `env:"PROJECT_ID" required:"true" default:""`
	PORT            string `env:"PORT" required:"false" default:"8080"`
}

func index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gofhir.WriteResponse(w, r, http.StatusOK, "Welcome to GoFHIR! You've hit the index page!")
	}
}
