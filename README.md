# e-commerce

This is a simple e-commerce service. It's api controller is built on top on gin-gonic and has two main route groups: orders and producs. Both support simple CRUD operations.

To fire up the server locally you will need to have docker configured on your machine. Once it's done, follow these steps:

1 - Bring database up:

```
make database
```
2 - Manually create database and collections: 
2.1 Database: order-service - Collection: orders
2.2 Database product-service - Collection: products

3 - Fire up the server:

```
make server
```

There is a built-in cadence workflow to handle a standard order workflow. The order workflow simulates a order lifecicle within a standard e-comerce service. It is designed to handle storage check and reservation, payment processing and order shipment.
