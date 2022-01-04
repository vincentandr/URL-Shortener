package paymentdb

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vincentandr/shopping-microservice/internal/model"
	"github.com/vincentandr/shopping-microservice/internal/mongodb"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
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
	paymentCollection := conn.Db.Collection("orders")

    return &Action{Conn: conn.Conn, Db: conn.Db, Collection: paymentCollection}, nil
}

func (a *Action) InitCollection(ctx context.Context) error {
	err := CreateIndex(ctx, a)
	if err != nil {
		return err
	}

	return nil
}

func CreateIndex(ctx context.Context, a *Action) error {
	// Create index
	_, err := a.Collection.Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{
                Keys: bson.D{
                        {Key: "user_id", Value: "text"},
                },
        },
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}

	return nil
}

func (a *Action) GetOrders(ctx context.Context, userId string) (*mongo.Cursor, error) {
	filter := bson.D{}
	if userId != "" {
		filter = bson.D{{Key: "user_id", Value: userId}}
	}

	res, err := a.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %v", err)
	}

	return res, nil
}

func (a *Action) GetItemsFromOrder(ctx context.Context, orderId primitive.ObjectID) (model.UserOrder, error) {
	var order model.UserOrder

	projection := bson.D{
		{Key:"user_id", Value: 1},
		{Key:"items", Value: 1},
	}

	if err := a.Collection.FindOne(ctx, bson.M{"_id":orderId}, options.FindOne().SetProjection(projection)).Decode(&order); err != nil{
		return model.UserOrder{}, fmt.Errorf("failed to get order: %v", err)
	}

	return order, nil
}

// Create order draft
func (a *Action) Checkout(ctx context.Context, userId string, items []*pb.ItemResponse) (string, error){
	var itemsBson []bson.M

	for _, item := range items {
		id, err := primitive.ObjectIDFromHex(item.ProductId)
		if err != nil{
			return "", fmt.Errorf("failed to convert from hex to objectID: %v", err)
		}

		itemsBson = append(itemsBson, bson.M{
			"_id": id,
		 	"name":item.Name,
			"price":item.Price, 
			"qty":item.Qty,
		})
	}

	// Create transaction function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Delete existing order draft document by userId
		_, err := a.Collection.DeleteOne(ctx, bson.M{"user_id": userId, "status": "draft"})
		if err != nil {
			return nil, fmt.Errorf("deletion of existing order draft failed: %v", err)
		}

		// Insert order document
		res, err := a.Collection.InsertOne(ctx, bson.M{
			"user_id": userId,
			"items": itemsBson,
			"status": "draft",
		})
		if err != nil {
			return nil, fmt.Errorf("order draft creation failed: %v", err)
		}

		return res.InsertedID, nil
	}
	// Start a transaction session
	session, err := a.Conn.StartSession()
	if err != nil {
		return "", err
	}
	defer session.EndSession(ctx)

	res, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return "", err
	}

	orderId := res.(primitive.ObjectID).Hex()

	return fmt.Sprintf("%v", orderId), nil
}

// Change status to paid
func (a *Action) MakePayment(ctx context.Context, orderId primitive.ObjectID) (error){
	_, err := a.Collection.UpdateOne(ctx, bson.M{"_id":orderId}, bson.M{"$set": bson.M{"status":"paid"}})
	if err != nil {
		return fmt.Errorf("failed to change order status: %v", err)
	}

	return nil
}