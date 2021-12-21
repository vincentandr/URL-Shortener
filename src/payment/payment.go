package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)


var schema = `
CREATE TABLE orders (
    id int NOT NULL AUTO_INCREMENT,
    user_id varchar(50),
	product_id varchar(50),
    qty int,
	total double,
	PRIMARY KEY (id)
);`

type Order struct {
	Order_Id int 
	User_id string 
	Product_id int
	Qty int
	Total float32
}

type mysqldb struct {
	db *sqlx.DB
}

func NewDb() (*mysqldb, error) {
	// User:pass@(addr:port)/database_name
	db, err := sqlx.Connect("mysql", "root@(127.0.0.1:3306)/product")
	if err != nil{
		return nil, err
	}

	return &mysqldb{db}, nil
}

func (m *mysqldb) initSchema() error {
	m.db.MustExec(schema)
	err := seedTable(m)
	
	return err
}

func (m *mysqldb) getAllOrders() ([]Order, error){
	orders := []Order{}

	err := m.db.Select(&orders, "select * from orders")

	return orders, err
}

func seedTable(m *mysqldb) error{
	tx := m.db.MustBegin()
	tx.MustExec("INSERT INTO orders (user_id, product_id, qty, total) VALUES (?, ?, ?, ?)", "user1", "1", 15, 4000)
    tx.MustExec("INSERT INTO orders (user_id, product_id, qty, total) VALUES (?, ?, ?, ?)", "user2", "2", 3, 500)
    err := tx.Commit()

	return err
}

func main() {
	// User:pass@(addr:port)/database_name
	db, err := NewDb()
    if err != nil {
        panic(err)
    }

	err = db.initSchema()
	if err != nil {
		fmt.Println("failed to create schema")
	}
}