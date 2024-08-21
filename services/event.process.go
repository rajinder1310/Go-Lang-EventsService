package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func (s *EventService) insertBulkData(Events []interface{}, collectionName *mongo.Collection) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	_, err := collectionName.InsertMany(ctx, Events, options.InsertMany().SetOrdered(false))
// 	if err != nil {
// 		log.Printf("Error inserting data: %v", err)
// 		return err
// 	}
// 	log.Println("Successfully inserted data.")
// 	return nil
// }

func (s *EventService) insertBulkData(events []interface{}, collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use ordered: false to continue inserting even if some documents fail
	_, err := collection.InsertMany(ctx, events, options.InsertMany().SetOrdered(false))
	if err != nil {
		// Check if the error is a bulk write error
		if bulkWriteErr, ok := err.(mongo.BulkWriteException); ok {
			// Log each write error
			for _, writeError := range bulkWriteErr.WriteErrors {
				log.Printf("Write error at index %d: %v", writeError.Index, writeError.Message)
			}
		} else {
			// Log any other errors
			log.Printf("Error inserting data: %v", err)
			return err
		}
	}

	log.Println("Successfully inserted data.")
	return nil
}
