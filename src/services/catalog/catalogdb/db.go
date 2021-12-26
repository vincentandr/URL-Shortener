package catalogdb

import (
	"context"
	"fmt"

	"github.com/vincentandr/shopping-microservice/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbConn *mongo.Client
	catalogDb *mongo.Database
	catalogCollection *mongo.Collection
)

func NewDb(ctx context.Context){
	// User:pass@(addr:port)/database_name
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:30001,localhost:30002,localhost:30003/?replicaSet=shop-mongo-set"))
	if err != nil{
		fmt.Print("failed to establish a new database connection")
		panic(err)
	}

	dbConn = conn
	
	// Create database
	catalogDb = dbConn.Database("catalog")
	catalogCollection = catalogDb.Collection("catalog")

	// Create index
	_, err = catalogCollection.Indexes().CreateOne(
        ctx,
        mongo.IndexModel{
                Keys: bson.D{{
                        Key:"name", Value: "text",
                }},
        },
	)
	if err != nil {
		fmt.Println("failed to create index")
	}

	// Seed collection
	if err = SeedCollection(ctx); err != nil {
		fmt.Println(err)
	}
}

func SeedCollection(ctx context.Context) error {
	count, err := catalogCollection.CountDocuments(ctx, bson.D{})
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

		_, err = catalogCollection.InsertMany(ctx, docs)
		if err != nil{
			return fmt.Errorf("failed to seed documents: %v", err)
		}
	} else if err != nil{
		return fmt.Errorf("failed to count documents: %v", err)
	}
	return nil
}

func Disconnect(ctx context.Context) {
	// Defer close connection
	if err := dbConn.Disconnect(ctx); err != nil{
		fmt.Print("failed to disconnect db connection: ")
		panic(err)
	}
}

func GetProducts(ctx context.Context) (*mongo.Cursor, error){
	cursor, err := catalogCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func GetProductsByIds(ctx context.Context, ids []primitive.ObjectID) (*mongo.Cursor, error){
	// Set what fields to get / select
	projection := bson.D{
		{Key:"_id", Value: 1},
		{Key:"name", Value: 1},
		{Key:"price", Value: 1},
	}

	cursor, err := catalogCollection.Find(
		ctx, 
		bson.M{"_id": bson.M{"$in": ids}},
		options.Find().SetProjection(projection),
	)
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func GetProductsByName(ctx context.Context, name string) (*mongo.Cursor, error){
	// equal to LIKE %name%
	cursor, err := catalogCollection.Find(ctx, bson.M{"name": primitive.Regex{Pattern: name, Options: ""}})
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func UpdateProducts(ctx context.Context, items []model.Product) (error){
	
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
		res, err := catalogCollection.BulkWrite(
			sessCtx,
			operations,
		)
		if err != nil {
			return nil, fmt.Errorf("order creation failed: %v", err)
		}

		if int(res.ModifiedCount) != len(items) {
			sessCtx.AbortTransaction(sessCtx)

			fmt.Println(int(res.ModifiedCount), len(items))
			return nil, fmt.Errorf("not all items qty are updated, not atomic")
		}

		return nil, nil
	}
	// Start a transaction session
	session, err := dbConn.StartSession()
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