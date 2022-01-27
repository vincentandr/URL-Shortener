package catalogdb

import (
	"context"
	_ "embed"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/vincentandr/shopping-microservice/internal/model"
	"github.com/vincentandr/shopping-microservice/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// embed catalog.csv into compiled go binary
//go:embed seed/catalog.csv
var catalog string

type Repository struct {
	Conn *mongo.Client
	Db *mongo.Database
	Collection *mongo.Collection
}

func NewRepository(conn *mongodb.Mongo) (*Repository, error) {
	catalogCollection := conn.Db.Collection("catalogs")

    return &Repository{Conn: conn.Conn, Db: conn.Db, Collection: catalogCollection}, nil
}

func (r *Repository) InitCollection(ctx context.Context) error {
	err := CreateIndex(ctx, r)
	if err != nil {
		return err
	}
	err = SeedCollection(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

func CreateIndex(ctx context.Context, r *Repository) error {
	// Create index
	_, err := r.Collection.Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{
                Keys: bson.D{{
                        Key:"name", Value: "text",
                }},
        },
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}

	return nil
}

func SeedCollection(ctx context.Context, r *Repository) error {
	count, err := r.Collection.CountDocuments(ctx, bson.D{})
	// Seed only if collection has 0 document
	if err == nil && count == 0 {
		docs := []interface{}{}

		// new strings reader from embedded catalog.csv
		stringsReader := strings.NewReader(catalog)
		csvReader := csv.NewReader(stringsReader)

		// read per line until EOF
		for{
			record, err := csvReader.Read()
			if err != nil { // EOF
				break
			}

			name := record[0]
			price, _ := strconv.Atoi(record[1])
			qty, _ := strconv.Atoi(record[2])
			desc := record[3]
			image := record[4]

			doc := bson.D{
				{Key:"name", Value: name},
				{Key:"price", Value: price},
				{Key:"qty", Value: qty},
				{Key:"desc", Value: desc},
				{Key:"image", Value: image},
			}
			docs = append(docs, doc)
		}

		_, err = r.Collection.InsertMany(ctx, docs)
		if err != nil{
			return fmt.Errorf("failed to seed documents: %v", err)
		}
	} else if err != nil{
		return fmt.Errorf("failed to count documents: %v", err)
	}
	return nil
}

func (r *Repository) GetProducts(ctx context.Context) (*mongo.Cursor, error){
	cursor, err := r.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func (r *Repository) GetProductsByIds(ctx context.Context, ids []primitive.ObjectID) (*mongo.Cursor, error){
	// Set what fields to get / select
	projection := bson.D{
		{Key:"_id", Value: 1},
		{Key:"name", Value: 1},
		{Key:"price", Value: 1},
	}

	cursor, err := r.Collection.Find(
		ctx, 
		bson.M{"_id": bson.M{"$in": ids}},
		options.Find().SetProjection(projection),
	)
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func (r *Repository) GetProductsByName(ctx context.Context, name string) (*mongo.Cursor, error){
	// equal to LIKE %name%
	cursor, err := r.Collection.Find(ctx, bson.M{"name": primitive.Regex{Pattern: name, Options: ""}})
	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return cursor, nil
}

func (r *Repository) UpdateProducts(ctx context.Context, items []model.Product) (error){
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
		res, err := r.Collection.BulkWrite(
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
	session, err := r.Conn.StartSession()
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