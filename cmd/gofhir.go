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
	defer fsClient.Close()

	// Instantiate interceptors
	demoInt := &interceptor.DemoInterceptor{
		RequestTime:     map[string]time.Time{},
		RequestDuration: map[string]int64{},
		DbQueryTime:     map[string]time.Time{},
		DbQueryDuration: map[string]int64{},
	}

	// Instantiate services and register interceptors
	p := patient.New(r, fsClient, APPLICATION_URL+":"+PORT, "patients")
	p.RegisterInterceptors(demoInt)

	// clClient, err := gofhir.NewCloudLoggerClient(ctx, PROJECT_ID)
	// if err != nil {
	// 	log.Fatalf("Failed to create client: %v", err)
	// 	return
	// }
	// defer clClient.Close()
	// logger := clClient.Logger("default").StandardLogger(logging.Info)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	data, _ := ioutil.ReadAll(r.Body)
	// 	defer r.Body.Close()
	// 	err = kafkaWriter.WriteMessages(ctx,
	// 		kafka.Message{
	// 			Key:   []byte(""),
	// 			Value: data,
	// 		},
	// 	)
	// 	if err != nil {
	// 		log.Fatalf("failed to write request to kafka: %v", err)
	// 	}
	// })

	log.Println("Starting server on port: " + PORT)
	http.ListenAndServe(":"+PORT, r)
}

func LoadConfigFromEnvironment() {
	godotenv.Load(".env")
	APPLICATION_URL = os.Getenv("APPLICATION_URL")
	PROJECT_ID = os.Getenv("PROJECT_ID")
	PORT = os.Getenv("PORT")

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
