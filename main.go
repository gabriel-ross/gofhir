package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
)

var PROJECT_ID = "gofhir"
var PORT = "8080"

func main() {
	r := chi.NewRouter()
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	http.ListenAndServe(":"+PORT, r)
}
