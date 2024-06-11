# Schwarz Challenge
### Made by <3 Juan Calcagno AKA Nacho

---

### What I've done
I've built (for sure is not fully completed. But it doens't contain as many bugs as what you guys showed me LOL) a "small" project where you guys can play with Shopping Cart and coupons.

### Packages I've used :detective:
	github.com/google/uuid v1.6.0
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/stretchr/testify v1.9.0
	go.uber.org/mock v0.4.0
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.10

### Postman Collection available :white_check_mark:
It is available in the repo.

### How to run it :scream_cat:
1. Have docker in your machine
2. `git clone` this repo
3. Once you are inside the repo
4. Run `docker-compose up -d` this will initiate a container with a running postgres DB
5. Run `make mod` so you download needed pkgs
6. Run `make migration-run dir=up` this will run all needed migrations
7. Run `make run` and if all good. Project should be running ready to get some HTTP calls.

### You don't want to run it? :smiling_imp:
1. Have docker in your machine
2. `git clone` this repo
3. Once you are inside the repo
4 Run `make test` and this will trigger a docker compose file that will spin up a test DB + mgirations and then run all the needed tests. By the time of writing this test are passing lol. 

### Structure :palm_tree:
```
📦schwarz-challenge
 ┣ 📂cmd
 ┃ ┗ 📂server
 ┃ ┃ ┣ 📜dev.go
 ┃ ┃ ┗ 📜main.go
 ┣ 📂internal
 ┃ ┣ 📂app
 ┃ ┃ ┣ 📜app.go
 ┃ ┃ ┣ 📜domain.go
 ┃ ┃ ┣ 📜http.go
 ┃ ┃ ┣ 📜infra.go
 ┃ ┃ ┗ 📜options.go
 ┃ ┣ 📂coupon
 ┃ ┃ ┣ 📜coupon.go
 ┃ ┃ ┗ 📜coupon_test.go
 ┃ ┣ 📂errors
 ┃ ┃ ┗ 📜errors.go
 ┃ ┣ 📂helpers
 ┃ ┃ ┗ 📜db.go
 ┃ ┣ 📂http
 ┃ ┃ ┣ 📜coupon.go
 ┃ ┃ ┣ 📜coupon_test.go
 ┃ ┃ ┣ 📜server.go
 ┃ ┃ ┣ 📜shopping_cart.go
 ┃ ┃ ┗ 📜shopping_cart_test.go
 ┃ ┣ 📂mocks
 ┃ ┃ ┣ 📜mock_coupon.go
 ┃ ┃ ┗ 📜mock_shopping_cart.go
 ┃ ┣ 📂postgres
 ┃ ┃ ┣ 📜db.go
 ┃ ┃ ┗ 📜db_test.go
 ┃ ┣ 📂repo
 ┃ ┃ ┣ 📜coupon.go
 ┃ ┃ ┣ 📜coupon_test.go
 ┃ ┃ ┣ 📜shopping_cart.go
 ┃ ┃ ┗ 📜shopping_cart_test.go
 ┃ ┣ 📂service
 ┃ ┃ ┣ 📜coupon.go
 ┃ ┃ ┣ 📜coupon_test.go
 ┃ ┃ ┣ 📜shopping_cart.go
 ┃ ┃ ┗ 📜shopping_cart_test.go
 ┃ ┗ 📂shopping_cart
 ┃ ┃ ┣ 📜shopping_cart.go
 ┃ ┃ ┗ 📜shopping_cart_test.go
 ┣ 📂migrations
 ┃ ┣ 📜20240608141014_init-svc.down.sql
 ┃ ┣ 📜20240608141014_init-svc.up.sql
 ┃ ┣ 📜20240608151129_add-main-tables.down.sql
 ┃ ┗ 📜20240608151129_add-main-tables.up.sql
 ┣ 📜.gitignore
 ┣ 📜Makefile
 ┣ 📜README.md
 ┣ 📜docker-compose.yml
 ┣ 📜docker-compose_test.yml
 ┣ 📜generate-mocks.sh
 ┣ 📜go.mod
 ┣ 📜go.sum
 ┗ 📜schwarz.postman_collection.json
 ```

### HTTP Endpoints :zap:
-  ***Shopping Cart***
```
// Creates a shopping cart 
POST localhost:8080/shopping-cart
```
Payload
```
{
    "items": [
        {
            "name": "item 1",
            "description": "test description",
            "price": 100
        },
        {
            "name": "item 2",
            "description": "test description",
            "price": 30.30
        }
    ]
}
``` 

```
// Returns a list of shopping carts
GET localhost:8080/shopping-cart
```

```
// Apply coupon to a shopping cart
PUT localhost:8080/shopping-cart/:id/apply-coupon/:coupon_id
```

---
- ***Coupon***
```
// Creates a coupon 
POST localhost:8080/coupon
```
Payload
```
{
    "name": "FREE30",
    "amount": 30
}
``` 

```
// Returns a list of coupons
GET localhost:8080/coupon
```
