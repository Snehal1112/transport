package client

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"testing"
)

func setup() *Connect{
	ctx := context.Background()
	url := "mongodb+srv://snehal:VXmWJuqjM8CdtzOa@erp-kod8k.mongodb.net/erp?retryWrites=true&w=majority"
	return NewConnection(ctx, url)
}

func TestNewConnection(t *testing.T) {
	conn := setup()
	err := conn.CreateDocument("leptop", bson.D{{"name", "sd"},{"name","ddd"}})
	if err != nil {
		log.Println("Error:", err)
	}

}

func TestConnect_DeleteByID(t *testing.T) {
	conn := setup()
	_, err := conn.DeleteByID("leptop", "5f120b174078a2f86e63fef6")
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_FindByID(t *testing.T) {
	conn := setup()
	_, err := conn.FindByID("leptop", "5f1270b5e45f4cb2bfa4ac19")
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_FindOne(t *testing.T) {
	conn := setup()
	_, err := conn.FindOne("leptop", bson.D{{"name","sd"}})
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_Search(t *testing.T) {
	conn := setup()
	_, err := conn.Search("leptop", bson.D{{"name","sd"}}, 0,0)
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_UpdateDocByID(t *testing.T) {
	conn := setup()
	result, err := conn.UpdateDocByID("leptop", "5f12fd8ecd9ee87cc6de167b", bson.D{{"name","snehal"}})
	if err != nil {
		log.Println("Error:", err)
	}

	log.Println(result)
}