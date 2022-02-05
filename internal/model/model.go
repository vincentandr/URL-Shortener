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
	Subtotal float32	`bson:"subtotal"`
	Customer Customer `bson:"customer"`
	Status      string                `bson:"status"`
}

type Customer struct {
	First_name string `bson:"first_name"`
	Last_name string `bson:"last_name"`
	Email string `bson:"email"`
	Address string `bson:"address"`
	Area string `bson:"area"`
	Postal string `bson:"postal"`
	Phone string `bson:"phone"`
}

type UserOrder struct {
	User_id     string             `bson:"user_id"`
	Items    []Product            `bson:"items"`
	Subtotal	float32	`bson:"subtotal"`
}