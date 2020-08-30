package client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
)

var (
	mgoClient    *mongo.Client
	databaseName = "erp"
)

// Connect used to initialize the connection.
type Connect struct {
	ctx context.Context
	url string

	log      *logrus.Logger
	loglevel string
}

// NewConnection function used to initialize the connection client.
func NewConnection(options ...Options) *Connect {
	conn := &Connect{}
	for _, option := range options {
		option(conn)
	}

	logLevel, err := logrus.ParseLevel(conn.loglevel)

	if err != nil {
		logLevel = logrus.ErrorLevel
	}

	conn.log = &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  "02/01/2006 03:04:12 PM",
		},
		Hooks: make(logrus.LevelHooks),
		Level: logLevel,
	}

	conn.log.WithFields(logrus.Fields{
		"url": conn.url,
	}).Info("Request to connect mongoDB server")

	return conn
}

func (conn *Connect) getMongoClient() (*mongo.Client, error) {
	if mgoClient == nil {
		client, err := mongo.Connect(conn.ctx, options.Client().ApplyURI(conn.url))
		if err != nil {
			conn.log.WithFields(logrus.Fields{
				"url": conn.url,
			}).Error("Unable to connect mongoDB server")
			return nil, err
		}
		mgoClient = client
	}
	return mgoClient, nil
}

// withDatabase function perform an action on given collection.
func (conn *Connect) withCollection(collection string, fn func(mongo.SessionContext, *mongo.Collection) error) error {
	client, err := conn.getMongoClient()
	if err != nil {
		return err
	}

	// Create a new Session and SessionContext.
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	session, err := client.StartSession(opts)
	if err != nil {
		conn.log.WithFields(logrus.Fields{
			"collection": collection,
		}).Info("Problem in creating session")
		return fmt.Errorf("error in creating new session: %w", err)
	}
	defer session.EndSession(conn.ctx)

	err = mongo.WithSession(conn.ctx, session, func(sessionCtx mongo.SessionContext) error {
		// TODO: insert some before and after hooks.
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
		opts := options.FindOne()
		return c.FindOne(ctx, q, opts).Decode(&result)
	}
	find := func() error {
		return conn.withCollection(collectionName, query)
	}

	if err = find(); err != nil {

		return nil, err
	}
	return result, nil
}

// FindByID function find the first matching document in collection based on given ID.
func (conn *Connect) FindByID(collectionName string, ID string) (result bson.M, err error) {
	query := func(ctx mongo.SessionContext, c *mongo.Collection) error {
		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return fmt.Errorf("invalid Object id %w", err)
		}
		return c.FindOne(ctx, bson.D{{"_id", objID}}).Decode(&result)
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

		if q == nil {
			q = bson.D{}
		}
		cursor, err := c.Find(ctx, q, fieldOptions)
		if err != nil {
			return err
		}

		for cursor.Next(ctx) {
			var document bson.M
			err = cursor.Decode(&document)
			if err != nil {
				log.Fatal(err)
			}
			searchResults = append(searchResults, document)
		}
		defer cursor.Close(ctx)
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

// DeleteByID function was delete the document record from
// given collection.
func (conn *Connect) DeleteByID(collectionName string, id string) (deleteCount int64, err error) {
	query := func(ctx mongo.SessionContext, c *mongo.Collection) error {
		var objID primitive.ObjectID
		objID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid Object id %w", err)
		}
		var result *mongo.DeleteResult
		result, err = c.DeleteOne(ctx, bson.D{{"_id", objID}}, options.Delete())
		deleteCount = result.DeletedCount
		return err
	}
	deleteAction := func() error {
		return conn.withCollection(collectionName, query)
	}
	if err = deleteAction(); err != nil {
		return 0, err
	}
	return deleteCount, nil
}

// UpdateDocByID function update the document by given id
func (conn *Connect) UpdateDocByID(collection string, id string, data interface{}) (updatedDoc bson.M, err error) {
	query := func(ctx mongo.SessionContext, c *mongo.Collection) error {
		var objID primitive.ObjectID
		objID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid Object id %w", err)
		}
		opts := options.FindOneAndUpdate().SetUpsert(true)
		err = c.FindOneAndUpdate(ctx, bson.D{{"_id", objID}}, bson.D{{"$set", data}}, opts).Decode(&updatedDoc)
		return err
	}

	updateAction := func() error {
		return conn.withCollection(collection, query)
	}

	if err = updateAction(); err != nil {
		return nil, err
	}

	return updatedDoc, nil
}
