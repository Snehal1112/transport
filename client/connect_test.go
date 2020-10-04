package client

import (
	"context"
	"log"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	DbURL  = ""
	DbName = ""
)

func init() {
	DbURL = os.Getenv("DB_URL")
	DbName = os.Getenv("DB_NAME")
}

func setup() *Connect {
	ctx := context.Background()

	if len(DbURL) == 0 {
		log.Fatal("DbName should not be empty")
	}

	if len(DbName) == 0 {
		log.Fatal("DbName should not be empty")
	}

	return NewConnection(
		WithDatabase(DbName),
		WithCtx(ctx),
		WithURL(DbURL),
		WithLogLevel("error"),
	)
}

func TestNewConnection(t *testing.T) {
	conn := setup()
	err := conn.CreateDocument("leptop", bson.D{{"name", "sd"}, {"name", "ddd"}})
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
	_, err := conn.FindOne("leptop", bson.D{{"name", "sd"}})
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_Search(t *testing.T) {
	conn := setup()
	data, err := conn.Search("leptop", nil, 0, 0)
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println(data)
}

func TestConnect_UpdateDocByID(t *testing.T) {
	conn := setup()
	result, err := conn.UpdateDocByID("leptop", "5f12fd8ecd9ee87cc6de167b", bson.D{{"name", "snehal"}})
	if err != nil {
		log.Println("Error:", err)
	}

	log.Println(result)
}
