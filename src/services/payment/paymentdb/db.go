package catalogdb

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vincentandr/shopping-microservice/src/model"
	pb "github.com/vincentandr/shopping-microservice/src/services/payment/paymentpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbConn *mongo.Client
	orderDb *mongo.Database
	ordersCollection *mongo.Collection
)

func NewDb(ctx context.Context){
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil{
		fmt.Print("failed to establish a new database connection")
		panic(err)
	}

	dbConn = conn
	
	// Create database
	orderDb = dbConn.Database("order")
	ordersCollection = orderDb.Collection("orders")

	// Create index
	_, err = ordersCollection.Indexes().CreateOne(
        ctx,
        mongo.IndexModel{
                Keys: bson.D{
                        {Key: "user_id", Value: "text"},
                },
        },
	)
	if err != nil {
		fmt.Println("failed to create index")
	}
}

func Disconnect(ctx context.Context) {
	if err := dbConn.Disconnect(ctx); err != nil{
		fmt.Print("failed to disconnect db connection: ")
		panic(err)
	}
}

func GetOrder(ctx context.Context, orderId primitive.ObjectID) (model.UserOrder, error) {
	var order model.UserOrder

	projection := bson.D{
		{Key:"user_id", Value: 1},
		{Key:"items", Value: 1},
	}

	if err := ordersCollection.FindOne(ctx, bson.M{"_id":orderId}, options.FindOne().SetProjection(projection)).Decode(&order); err != nil{
		return model.UserOrder{}, fmt.Errorf("failed to get order: %v", err)
	}

	return order, nil
}

// Create order draft
func Checkout(ctx context.Context, userId string, items []*pb.ItemResponse) (string, error){
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
		_, err := ordersCollection.DeleteOne(ctx, bson.M{"user_id": userId, "status": "draft"})
		if err != nil {
			return nil, fmt.Errorf("deletion of existing order draft failed: %v", err)
		}

		// Insert order document
		res, err := ordersCollection.InsertOne(ctx, bson.M{
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
	session, err := dbConn.StartSession()
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
func MakePayment(ctx context.Context, orderId string) (error){
	return nil
}