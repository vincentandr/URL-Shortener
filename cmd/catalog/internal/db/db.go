package catalogdb

import (
	"context"
	"fmt"

	"github.com/vincentandr/shopping-microservice/internal/model"
	"github.com/vincentandr/shopping-microservice/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Action struct {
	Conn *mongo.Client
	Db *mongo.Database
	Collection *mongo.Collection
}

func NewAction(conn *mongodb.Mongo) (*Action, error) {
	catalogCollection := conn.Db.Collection("catalogs")

	// Create index
	_, err := catalogCollection.Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{
                Keys: bson.D{{
                        Key:"name", Value: "text",
                }},
        },
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

    return &Action{Conn: conn.Conn, Db: conn.Db, Collection: catalogCollection}, nil
}

func (a *Action) SeedCollection(ctx context.Context) error {
	count, err := a.Collection.CountDocuments(ctx, bson.D{})
	if err == nil && count == 0 {
		docs := []interface{}{
			bson.D{
				{Key:"name", Value: "laptop"},
				{Key:"price", Value: 600},
				{Key:"qty", Value: 14},
			},
			bson.D{
				{Key:"name", Value: "computer"},
				{Key:"price", Value: 800},
				{Key:"qty", Value: 32},
			},
		}

		_, err = a.Collection.InsertMany(ctx, docs)
		if err != nil{
			return fmt.Errorf("failed to seed documents: %v", err)
		}
	} else if err != nil{
		return fmt.Errorf("failed to count documents: %v", err)
	}
	return nil
}

func (a *Action) GetProducts(ctx context.Context) (*mongo.Cursor, error){
	cursor, err := a.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func (a *Action) GetProductsByIds(ctx context.Context, ids []primitive.ObjectID) (*mongo.Cursor, error){
	// Set what fields to get / select
	projection := bson.D{
		{Key:"_id", Value: 1},
		{Key:"name", Value: 1},
		{Key:"price", Value: 1},
	}

	cursor, err := a.Collection.Find(
		ctx, 
		bson.M{"_id": bson.M{"$in": ids}},
		options.Find().SetProjection(projection),
	)
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func (a *Action) GetProductsByName(ctx context.Context, name string) (*mongo.Cursor, error){
	// equal to LIKE %name%
	cursor, err := a.Collection.Find(ctx, bson.M{"name": primitive.Regex{Pattern: name, Options: ""}})
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func (a *Action) UpdateProducts(ctx context.Context, items []model.Product) (error){
	var operations []mongo.WriteModel

	for _, item := range items {
		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(
			bson.D{
				{Key: "_id", Value: item.Product_id}, 
				{Key: "qty", Value: bson.M{"$gte": item.Qty}},
			},
		).SetUpdate(
			bson.D{{Key: "$inc", Value: bson.D{{Key: "qty", Value: item.Qty * -1}}}},
		)

		operations = append(operations, operation)
	}
	
	// Create transaction function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Create new update operation for each cart item
		res, err := a.Collection.BulkWrite(
			sessCtx,
			operations,
		)
		if err != nil {
			return nil, fmt.Errorf("order creation failed: %v", err)
		}

		if int(res.ModifiedCount) != len(items) {
			sessCtx.AbortTransaction(sessCtx)
		}

		return nil, nil
	}
	// Start a transaction session
	session, err := a.Conn.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return nil
}