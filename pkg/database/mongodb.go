package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	
	clientOptions.SetMaxPoolSize(50)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetSocketTimeout(30 * time.Second)
	
	clientOptions.SetRetryWrites(true)
	clientOptions.SetRetryReads(true)

	log.Printf("Connecting to MongoDB Atlas...")
	log.Printf("Database: %s", dbName)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB Atlas: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB Atlas: %v", err)
	}

	log.Printf("Successfully connected to MongoDB Atlas")

	database := client.Database(dbName)

	if err := createAtlasIndexes(ctx, database); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	if err := logAtlasInfo(ctx, client, dbName); err != nil {
		log.Printf("Warning: Failed to get Atlas info: %v", err)
	}

	return database, nil
}

func createAtlasIndexes(ctx context.Context, db *mongo.Database) error {
	userCollection := db.Collection("users")
	
	// Create unique index for email field
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("email_unique"),
	}

	// Create compound index for email and is_active (common query pattern)
	emailActiveIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}, {Key: "is_active", Value: 1}},
		Options: options.Index().SetName("email_active_compound"),
	}

	// Create index for created_at (for sorting and filtering)
	createdAtIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: -1}},
		Options: options.Index().SetName("created_at_desc"),
	}

	// Create index for is_active (for filtering active users)
	activeIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "is_active", Value: 1}},
		Options: options.Index().SetName("is_active"),
	}

	// Create text search index for name and email
	textSearchIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: "text"},
			{Key: "email", Value: "text"},
		},
		Options: options.Index().SetName("name_email_text_search"),
	}

	indexes := []mongo.IndexModel{
		emailIndex,
		emailActiveIndex,
		createdAtIndex,
		activeIndex,
		textSearchIndex,
	}

	// Create all indexes
	_, err := userCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create Atlas indexes: %v", err)
	}

	log.Println("MongoDB Atlas indexes created successfully")
	return nil
}

func logAtlasInfo(ctx context.Context, client *mongo.Client, dbName string) error {
	// Get server status
	var serverStatus bson.M
	err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "serverStatus", Value: 1}}).Decode(&serverStatus)
	if err != nil {
		return err
	}

	if host, ok := serverStatus["host"].(string); ok {
		log.Printf("Atlas Host: %s", host)
	}

	if version, ok := serverStatus["version"].(string); ok {
		log.Printf("MongoDB Version: %s", version)
	}

	// Get database stats
	var dbStats bson.M
	err = client.Database(dbName).RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&dbStats)
	if err != nil {
		return err
	}

	if collections, ok := dbStats["collections"].(int32); ok {
		log.Printf("Collections: %d", collections)
	}

	return nil
}

func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB Atlas: %v", err)
	}

	log.Println("Disconnected from MongoDB Atlas")
	return nil
}

func GetCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}

func HealthCheck(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.Client().Ping(ctx, readpref.Primary())
}

func GetConnectionStats(db *mongo.Database) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var serverStatus bson.M
	err := db.Client().Database("admin").RunCommand(ctx, bson.D{{Key: "serverStatus", Value: 1}}).Decode(&serverStatus)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})
	
	if connections, ok := serverStatus["connections"].(bson.M); ok {
		stats["connections"] = connections
	}

	if network, ok := serverStatus["network"].(bson.M); ok {
		stats["network"] = network
	}

	stats["timestamp"] = time.Now().UTC().Format("2006-01-02 15:04:05")
	
	return stats, nil
}