package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)


var schema = `
CREATE TABLE products (
    id int NOT NULL AUTO_INCREMENT,
    name varchar(50),
    qty text,
	PRIMARY KEY (id)
);`

type Product struct {
	Id int 
	Name string 
	Qty int 
}

func main() {
	// User:pass@(addr:port)/database_name
	db, err := sqlx.Connect("mysql", "root@(127.0.0.1:3306)/shop")
    if err != nil {
        log.Fatalln(err)
    }

	// db.MustExec(schema)

	products := []Product{}

	tx := db.MustBegin()
    tx.MustExec("INSERT INTO products (name, qty) VALUES (?, ?)", "laptop", 15)
    tx.MustExec("INSERT INTO products (name, qty) VALUES (?, ?)", "computer", 3)
    tx.Commit()

	db.Select(&products, "select * from products")

	log.Println("products...")
	log.Println(products)
}