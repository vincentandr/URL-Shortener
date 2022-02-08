package mongodb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Conn *mongo.Client
	Db *mongo.Database
}

func NewDb(ctx context.Context, database string) (*Mongo, error){
	// User:pass@(addr:port)/database_name
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil{
		return nil, fmt.Errorf("failed to establish a new database connection: %v", err)
	}
	
	// Create database
	db := conn.Database(database)

	return &Mongo{
		Conn: conn,
		Db: db,
	}, nil
}

func (m *Mongo) Close() error {
	// Defer close connection
	if err := m.Conn.Disconnect(context.Background()); err != nil{
		return fmt.Errorf("failed to disconnect db connection: ")
	}

	return nil
}