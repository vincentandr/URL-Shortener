# shopping-microservice (In Progress)

![Untitled Diagram](https://user-images.githubusercontent.com/42005057/147475016-bc2f7406-ed5d-4d1f-8471-b9ba47caef6b.png)

## API

### Catalog APIs
```
/products GET --- Get all products 
/products/search GET --- Search products by name query parameter
```

### Cart APIs
```
/cart/{userId} GET --- Get cart items
/cart/{userId} DELETE --- Delete all cart items
/cart/{userId}/{productId} PUT --- Add or replace cart item with qty query parameter
/cart/{userId}/{productId} DELETE --- Delete a particular cart item by id
/cart/checkout/{userId} GET --- Checkout, returns an order ID to be used for payment
```

### Payment APIs
```
/payment/{orderId} PUT --- Change order status to paid
```

## Proto compiler command

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <filename>
```

## Docker commands

### MongoDB
Create replica set for session transactions
```
docker network create shop-mongo-cluster
```

Run 3 containers
```
docker run -d --net my-mongo-cluster -p 27017:27017 --name mongo1 mongo:5.0.5-focal mongod --replSet shop-mongo-set --port 27017
docker run -d --net my-mongo-cluster -p 27018:27018 --name mongo2 mongo:5.0.5-focal mongod --replSet shop-mongo-set --port 27018
docker run -d --net my-mongo-cluster -p 27019:27019 --name mongo3 mongo:5.0.5-focal mongod --replSet shop-mongo-set --port 27019
```

Connect to mongo shell in one of the containers
```
docker exec -it mongo1 mongo
```

Initiate a replica set with the config
```
config = {
  	"_id" : "shop-mongo-set",
  	"members" : [
  		{
  			"_id" : 0,
  			"host" : "mongo1:27017"
  		},
  		{
  			"_id" : 1,
  			"host" : "mongo2:27018"
  		},
  		{
  			"_id" : 2,
  			"host" : "mongo3:27019"
  		}
  	]
  }
  
rs.initiate(config)
```

Note: these lines may need to be added in your OS `etc/hosts` host file
```
127.0.0.1 mongo1
127.0.0.1 mongo2
127.0.0.1 mongo3
```

### Redis
```
docker run --name cart-redis -d redis
```

### RabbitMQ
15672 for RabbitMQ Manager and 5672 for connection requests
```
docker run -d --hostname shop-rabbitmq --name shop-rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3.9.11-management
```

## Run main.go
Run these files from the root project directory

Catalog microservice
```
go run cmd/catalog/catalog.go
```

Cart microservice
```
go run cmd/cart/cart.go
```

Payment microservice
```
go run cmd/payment/payment.go
```

Backend for frontend microservice
```
go run cmd/bff/bff.go
```
