package catalogdb

import (
	"log"

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
	Qty int
}

func NewDb(){
	// User:pass@(addr:port)/database_name
	db, err := sqlx.Connect("mysql", "root@(127.0.0.1:3306)/product")
	if err != nil{
		panic("Failed to establish a database connection")
	}

	dbConn = db
}

func Disconnect() error{
	err := dbConn.Close()
	if err != nil{
		log.Fatalln(err)
	}

	return err
}

func InitSchema() error{
	dbConn.MustExec(schema)
	err := SeedTable()
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func SeedTable() error{
	tx := dbConn.MustBegin()
	tx.MustExec("INSERT INTO products (name, qty) VALUES (?, ?)", "laptop", 15)
    tx.MustExec("INSERT INTO products (name, qty) VALUES (?, ?)", "computer", 3)
    err := tx.Commit()
	
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetProducts() ([]Product, error){
	products := []Product{}

	err := dbConn.Select(&products, "select * from products")

	if err != nil {
		log.Fatalln("GetProducts Select query failed", err)
		return nil, err
	}

	return products, nil
}

func GetProductsWithName(name string) ([]Product, error){
	products := []Product{}

	err := dbConn.Select(&products, "select * from products where name like ?", "%"+name+"%")

	if err != nil {
		log.Fatalln("GetProductsWithName Select query failed", err)
		return nil, err
	}

	return products, nil
}