package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Product_id primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Price      float32            `bson:"price"`
	Qty        int                `bson:"qty"`
}