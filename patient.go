package gofhir

type Patient struct {
	ID   string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
}
