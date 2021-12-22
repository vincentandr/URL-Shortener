package catalogdb

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	dbConn *sqlx.DB
	schema string = `
CREATE TABLE products (
    product_id int NOT NULL AUTO_INCREMENT,
    name varchar(50),
    qty int,
	PRIMARY KEY (product_id)
);`
)

// Fields must have capital letter to be exported and used in another package 
type Product struct {
	Product_id int
	Name string
	Price float32
	Qty int
}

type ProductById struct {
	Product_id int
	Name string
	Price float32
}

func NewDb(){
	// User:pass@(addr:port)/database_name
	db, err := sqlx.Connect("mysql", "root@(127.0.0.1:3306)/product")
	if err != nil{
		fmt.Print("failed to establish a new database connection: ")
		panic(err)
	}

	dbConn = db

	// Defer close connection
	defer func(){
		if err := dbConn.Close(); err != nil{
			fmt.Print("failed to disconnect db connection: ")
			panic(err)
		}
	}()
}

func InitSchema() error{
	dbConn.MustExec(schema)
	err := SeedTable()
	if err != nil {
		return fmt.Errorf("create schema error: %v", err)
	}

	return nil
}

func SeedTable() error{
	tx := dbConn.MustBegin()
	tx.MustExec("INSERT INTO products (name, qty) VALUES (?, ?)", "laptop", 15)
    tx.MustExec("INSERT INTO products (name, qty) VALUES (?, ?)", "computer", 3)
    err := tx.Commit()
	
	if err != nil {
		return fmt.Errorf("seed table error: %v", err)
	}

	return nil
}

func GetProducts(ctx context.Context) ([]Product, error){
	products := []Product{}

	err := dbConn.SelectContext(ctx, &products, "select * from products")

	if err != nil {
		return nil, fmt.Errorf("GetProducts Select query failed: %v", err)
	}

	return products, nil
}

func GetProductsByIds(ctx context.Context, ids []int) ([]ProductById, error){
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

func GetProductsByName(ctx context.Context, name string) ([]Product, error){
	products := []Product{}

	err := dbConn.SelectContext(ctx, &products, "select * from products where name like ?", "%"+name+"%")

	if err != nil {
		return nil, fmt.Errorf("GetProductsByName Select query failed: %v", err)
	}

	return products, nil
}