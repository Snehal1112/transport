package client

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"testing"
)

func TestNewConnection(t *testing.T) {
	ctx := context.Background()

	url := "mongodb+srv://snehal:VXmWJuqjM8CdtzOa@erp-kod8k.mongodb.net/erp?retryWrites=true&w=majority"
	conn := NewConnection(ctx, url)
	err := conn.CreateDocument("leptop", bson.D{{"name", "Alice"}})
	if err != nil {
		log.Println("Error:", err)
	}

}