package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Product_id primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Price      float32            `bson:"price"`
	Qty        int                `bson:"qty"`
	Desc string `bson:"desc"`
	Image string `bson:"image"`
}

type Order struct {
	Order_id primitive.ObjectID `bson:"_id"`
	User_id     string             `bson:"user_id"`
	Items    []Product            `bson:"items"`
	Status      string                `bson:"status"`
}

type UserOrder struct {
	User_id     string             `bson:"user_id"`
	Items    []Product            `bson:"items"`
}