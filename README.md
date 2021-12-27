# shopping-microservice (In Progress)

![Untitled Diagram drawio(1)](https://user-images.githubusercontent.com/42005057/147434228-5b1902c0-da80-42dd-a468-ad3f998c3f23.png)

## Proto compiler command

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <filename>`

## Docker commands

### MySQL
`docker run --name shopdb -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -e MYSQL_DATABASE=product -p 3306:3306 -d mysql:8.0.27`

### PHPMyAdmin
`docker run --name myshopadmin --link shopdb:mysql -p 80:80 -e PMA_HOST=shopdb phpmyadmin:5.1.1-apache`
