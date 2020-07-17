package client

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mgoClient    *mongo.Client
	databaseName = "erp"
)

// Connect used to initialize the connection.
type Connect struct {
	ctx context.Context
	url string
}

// NewConnection function used to initialize the connection client.
func NewConnection(ctx context.Context, url string) *Connect {
	return &Connect{ctx, url}
}

func (conn *Connect) getMongoClient() *mongo.Client {
	if mgoClient == nil {
		client, err := mongo.Connect(conn.ctx, options.Client().ApplyURI(conn.url))
		if err != nil {
			log.Fatal(err)
		}
		mgoClient = client
	}

	return mgoClient
}

// withDatabase function perform an action on given collection.
func (conn *Connect) withCollection(collection string, fn func(mongo.SessionContext, *mongo.Collection) error) error {
	client := conn.getMongoClient()

	// Create a new Session and SessionContext.
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	session, err := client.StartSession(opts)
	if err != nil {
		err := fmt.Errorf("Error in creating new session: %w", err)
		log.Println(err)
	}
	ctx := conn.ctx
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionCtx mongo.SessionContext) error {
		c := client.Database(databaseName).Collection(collection)
		return fn(sessionCtx, c)
	})
	return err
}

// CreateDocument function was create the new document record in
// given collection.
func (conn *Connect) CreateDocument(collectionName string, props interface{}) error {
	query := func(ctx mongo.SessionContext, c *mongo.Collection) error {
		_, err := c.InsertOne(ctx, props)
		return err
	}
	insert := func() error {
		return conn.withCollection(collectionName, query)
	}
	return insert()
}

// FindOne function find the first matching document in collection based on given query.
func (conn *Connect) FindOne(collectionName string, q interface{}) (result bson.M, err error) {
	query := func(ctx mongo.SessionContext, c *mongo.Collection) error {
		return c.FindOne(ctx, q).Decode(result)
	}
	find := func() error {
		return conn.withCollection(collectionName, query)
	}

	if err = find(); err != nil {
		return nil, err
	}
	return result, nil
}

// Search function search the document based on given search query.
func (conn *Connect) Search(collectionName string, q interface{}, skip int, limit int) (searchResults []bson.M, err error) {
	query := func(ctx mongo.SessionContext, c *mongo.Collection) error {
		fieldOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
		cursor, err := c.Find(ctx, q, fieldOptions)
		_ = cursor.Decode(searchResults)
		return err
	}

	search := func() error {
		return conn.withCollection(collectionName, query)
	}

	if err = search(); err != nil {
		return nil, err
	}

	return searchResults, nil
}