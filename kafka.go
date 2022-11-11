package gofhir

// kafkaWriter := &kafka.Writer{
// 	Addr:                   kafka.TCP("localhost:9092"),
// 	Topic:                  "requests",
// 	AllowAutoTopicCreation: true,
// 	Balancer:               &kafka.LeastBytes{},
// }

// err = kafkaWriter.WriteMessages(ctx,
// 	kafka.Message{
// 		Key:   []byte("some key"),
// 		Value: []byte("some value"),
// 	},
// )
// if err != nil {
// 	log.Fatalf("failed to write messages: %v", err)
// 	return
// }

// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 	err = kafkaWriter.WriteMessages(ctx,
// 		kafka.Message{
// 			Key:   []byte(""),
// 			Value: []byte(r.RemoteAddr),
// 		},
// 	)
// 	if err != nil {
// 		log.Fatalf("failed to write request to kafka: %v", err)
// 	}
// })
