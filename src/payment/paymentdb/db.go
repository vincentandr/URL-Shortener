package catalogdb

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	pb "github.com/vincentandr/shopping-microservice/src/payment/paymentpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbConn *mongo.Client
	orderDb *mongo.Database
	ordersCollection *mongo.Collection
	orderSchema string = `
	CREATE TABLE orders (
		order_id int NOT NULL AUTO_INCREMENT,
		user_id varchar(50),
		timestamp date,
		PRIMARY KEY (order_id)
	);`
)

// Fields must have capital letter to be exported and used in another package 
type Item struct {
	ProductId int
	Name string
	Price float32
	Qty int
}

func NewDb(){
	// User:pass@(addr:port)/database_name
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()	

	conn, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil{
		fmt.Print("failed to establish a new database connection: %v", err)
		panic(err)
	}

	dbConn = conn

	// Defer close connection
	defer func(){
		if err := dbConn.Disconnect(ctx); err != nil{
			fmt.Print("failed to disconnect db connection: ")
			panic(err)
		}
	}()
	
	// Create database
	orderDb = dbConn.Database("order")
	ordersCollection = orderDb.Collection("orders")
}

// Create order draft
func Checkout(ctx context.Context, userId string, items []*pb.ItemResponse) (string, error){
	var itemsBson []bson.M

	for _, item := range items {
		itemsBson = append(itemsBson, bson.M{"product_id":item.ProductId, "name":item.Name, "price":item.Price, "qty":item.Qty})
	}

	// Generate order ID
	orderId := uuid.New().String()

	// Insert document
	_, err := ordersCollection.InsertOne(ctx, bson.D{
		{"order_id", orderId},
		{"user_id", userId},
		{"items", itemsBson},
		{"status", "draft"},
	})
	if err != nil {
		return "", fmt.Errorf("order draft creation failed: %v", err)
	}

	return orderId, nil
}

// Change status to paid
func MakePayment(ctx context.Context, orderId string) (error){
	products := []ProductById{}

	query, args, err := sqlx.In("select product_id, name from products where product_id in (?)", ids)
	if err != nil {
		return nil, fmt.Errorf("select IN clause error: %v", err)
	}

	err = dbConn.SelectContext(ctx, &products, dbConn.Rebind(query), args...)

	if err != nil {
		return nil, fmt.Errorf("GetProductsByIds Select query failed: %v", err)
	}

	return products, nil
}