package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir/patient"
	"github.com/gabriel-ross/gofhir/storage"
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

	pc := storage.NewPatientClient(client)
	pch := storage.NewHooker(*pc)
	pch.RegisterBeforeDBTransactionInterceptor(HelloWorld)
	svc := patient.New(pch)
	r.Mount("/patients", svc.Routes())

	http.ListenAndServe(":"+PORT, r)
}

func HelloWorld() {
	log.Println("hello world")
}
