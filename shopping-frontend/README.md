
# shopping-microservice (In Progress)

## Frontend UI

https://github.com/vincentandr/shopping-frontend

## Tech Stack Diagram
![Untitled Diagram](https://user-images.githubusercontent.com/42005057/147475016-bc2f7406-ed5d-4d1f-8471-b9ba47caef6b.png)

## Start via Docker commands
**Create `data` folder for mongo DB replica set first:**
```
shopping-microservice/
  ...
  data/
    mongo1/
    mongo2/
    mongo3/
```
Then start the docker containers by using the command:
```
docker-compose up -d
```

## API

Call the APIs using API platform, e.g. Postman.

### Catalog APIs
```
/products                   GET - Get all products 
/products/search            GET - Search products by name query parameter
```

### Cart APIs
```
/cart/{userId}              GET - Get cart items
/cart/{userId}              DELETE - Delete all cart items
/cart/{userId}/{productId}  PUT - Add or replace cart item with qty query parameter
/cart/{userId}/{productId}  DELETE - Delete a particular cart item by id
/cart/checkout/{userId}     GET - Checkout, returns an order ID to be used for payment
```

### Payment APIs
```
/payment                    GET - Get all orders
/payment/{userId}           GET - Get all orders by user ID
/payment/draft/{userId}     GET - Get the user's draft order. Beginning of checkout, and creating Stripe payment flow
/payment/{orderId}          POST - Make payment
```

## Tests
HTTP handlers test
```
go test -v ./cmd/bff/internal/handler/...
```

Cart helper functions test
```
go test -v ./cmd/cart/
```

## Miscellaneous
### RabbitMQ Management / Admin

Access RabbitMQ admin web page for monitoring purpose.
```
localhost:15672
```

### Proto compiler command

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <filename>
```

## Possible improvements

- Purely asynchronous communication between microservices (RabbitMQ pub/sub instead of GRPC, and return the async response using websocket / the client periodically polls for result)
- Caching results for certain requests

# Shopping-frontend
Frontend for [shopping-microservice](https://github.com/vincentandr/shopping-microservice) system done in React, Redux, and Material UI.

## Screenshots

1. Homepage

![image](https://user-images.githubusercontent.com/42005057/152652903-1dc97811-cac2-475c-8cf0-f7ec4b2f6942.png) 

![image](https://user-images.githubusercontent.com/42005057/152685650-89c80910-f908-412c-8306-dca7adfb70e0.png)

2. Cart Drawer open

![image](https://user-images.githubusercontent.com/42005057/152652929-a4e6d16f-8340-4149-82a5-9a1bdd4a67aa.png)

![image](https://user-images.githubusercontent.com/42005057/152685695-297f79b0-7788-4501-b0c8-e2fabbbf1a4e.png)

3. Checkout Form 1

![image](https://user-images.githubusercontent.com/42005057/152685740-806db05f-5fc1-448e-84bc-bb6460ff63de.png)

![image](https://user-images.githubusercontent.com/42005057/152685723-6b9b3f1c-f470-41dd-ba11-4ebdea39d1b2.png)


4. Checkout Form 2 (Stripe payment)


![image](https://user-images.githubusercontent.com/42005057/152685762-aee56bc5-52d8-480a-a3f6-704e0e5c4f25.png)
