# shopping-microservice

## Docker commands

### MySQL
`docker run --name db -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -e MYSQL_DATABASE=shop -p 3306:3306 -d mysql:latest`

### PHPMyAdmin
`docker run --name myadmin --link db:mysql -p 80:80 -e PMA_HOST=db phpmyadmin:latest`
