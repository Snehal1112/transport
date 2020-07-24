package client

import (
	"context"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func setup() *Connect {
	ctx := context.Background()
	viper.AddConfigPath("")
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	url := "mongodb+srv://snehal:VXmWJuqjM8CdtzOa@erp-kod8k.mongodb.net/erp?retryWrites=true&w=majority"
	return NewConnection(ctx, url)
}

func insertMockData() {
	conn := setup()
	var objID primitive.ObjectID
	objID, _ = primitive.ObjectIDFromHex("5f120b174078a2f86e63fef6")
	conn.DeleteByID("leptop", "5f120b174078a2f86e63fef6")
	err := conn.CreateDocument("leptop", bson.D{{"_id", objID},{"first_name", "snehal"}, {"last_name", "dangroshiya"}})
	if err != nil {
		log.Println("Error:", err)
	}

}

func TestNewConnection(t *testing.T) {
	conn := setup()
	err := conn.CreateDocument("leptop", bson.D{{"first_name", "data1"}, {"last_name", "data2"}})
	if err != nil {
		log.Println("Error:", err)
	}

}

func TestConnect_DeleteByID(t *testing.T) {
	conn := setup()
	insertMockData()
	_, err := conn.DeleteByID("leptop", "5f120b174078a2f86e63fef6")
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_FindByID(t *testing.T) {
	conn := setup()
	insertMockData()
	_, err := conn.FindByID("leptop", "5f120b174078a2f86e63fef6")
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_FindOne(t *testing.T) {
	conn := setup()
	_, err := conn.FindOne("leptop", bson.D{{"first_name", "snehal"}})
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_Search(t *testing.T) {
	conn := setup()
	_, err := conn.Search("leptop", bson.D{{"name", "sd"}}, 0, 0)
	if err != nil {
		log.Println("Error:", err)
	}
}

func TestConnect_UpdateDocByID(t *testing.T) {
	conn := setup()
	insertMockData()
	doc , err := conn.UpdateDocByID("leptop", "5f120b174078a2f86e63fef6", bson.D{{"first_name", "dhara"}})
	if err != nil {
		log.Println("Error:", err)
	}
	result, _ :=bson.MarshalExtJSON(&doc, true, true)

	objID, _ := primitive.ObjectIDFromHex("5f120b174078a2f86e63fef6")
	mockData, _ :=bson.MarshalExtJSON(&bson.D{{"_id", objID},{"first_name", "ssnehal"}, {"last_name", "dangroshiya"}}, true, true)

	assert.Equal(t,string(result),  string(mockData))
}
