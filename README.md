# Solve Race Conditions in OrderManagementSystem
-------
Problem: 

An eCommerce company is having a sale, where some products are having lightning deals with an extra discount for 30 mins. When a product is having a lightning deal, it is listed in the landing page and a huge amount of traffic is bombarded on the product’s detail page. One such product X is having stock of 10 and 1,00,000 users have landed on its product details page out of which at time zero 10,000 users have added product in their cart and proceed to checkout. Since available quality is only 10, there will be a race condition among 10,000 users to get 10 quantities. 

Below solution solve the race condition, can handle high concurrency with ease,microservices talk to each other


-------

It has simple dependencies:

 - [GRPC library )](google.golang.org/grpc)

Get Started:


-------

Clone the source

    https://github.com/vinay10949/OrderManagementSystem.git

Setup dependencies

    go get google.golang.org/grpc
   
Run the app

    go run checkProductAvailableServer.go 
    go run addToCartAvailableServer.go 
    go run purchaseProductServer.go 
    go run client.go 


----------

[Folder Structure](https://irahardianto.github.io/service-pattern-go/#folder-structure)
-------
    /
    |- addToCart
    |- availablequantity
    |- purchase
    checkProductAvailableServer.go
    addToCartAvailableServer.go
    purchaseProductServer.go
    client.go


Every folder is a namespace of their own, and every file / struct under the same folder should only use the same namepace as their root folder.

### addToCart

This has proto request and interface methods for adding product to cart


### availablequantity

This has proto request and interface methods for checking availablity of item in stock

### purchase

This has proto request and interface methods for purchase of item from cart 

### checkProductAvailableServer.go

This is productAvailablity server,that manages inventory of products, it will say if product is available or not ,if available after purchase according stock will be updated

### addToCartAvailableServer.go

This is addToCart server,that allows user to add product to cart

### purchase.go

this will allows user to purchase the product 

### client.go

This is our client which interacts with these service


###Interaction

client---callsOrderProcess---->callsAddTocart-->Purchase

Service addToCart Interacts with service checkProductAvailablity
Service purchase Interacts with service addToCart and then interacts with checkProductAvailablity to update no of items left in stock



----------
