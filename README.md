
# e-commerce
## Server
This is a simple e-commerce service. It's api controller is built on top on gin-gonic and has two main route groups: orders and products. Both support simple CRUD operations.

To fire up the server locally you will need to have docker configured on your machine. Once it's done, follow these steps:

1 - Bring database up:

```
make database
```

2 - Fire up the server:

```
make server
```

3 - Create a product by sending a POST reques to to http://localhost:8080/products with the following body:
```json
{
	"productId":  "dbc9753d-2b2a-487c-a6bc-c8bffd901c68",
	"name":  "iPhone 13 Special Moon Edition",
	"units":  72,
	"price":  9000
}
```

4 - Create a order by sending a POST reques to to http://localhost:8080/orders with the following body:
```json
{
	"orderId":  "c52c657f-dbb6-4027-8844-9e1f4c97fe59",
	"userId":  "c0b98628-2509-4d3d-a8b2-b8c40147b682",
	"productId":  "dbc9753d-2b2a-487c-a6bc-c8bffd901c68",
	"status":  "pending",
	"quantity":  2,
	"deliveryAddress":  "Av. Pres. Antônio Carlos, 6627 - Pampulha, Belo Horizonte - MG, 31270-901"
}
```



## Cadence
There is a built-in cadence workflow to handle a standard order workflow. The order workflow simulates a order lifecycle within a standard e-comerce service. It is designed to handle storage check and reservation, payment processing and order shipment.

To fire up a order workflow locally you will need to have docker configured on your machine. Once it's done, follow these steps:

1 - Bring database up:

```
make database
```

2 - Start all cadence containers needed to run the workflow:

```
make cadence-containers
```

3 - Register the order-service domain in cadence server:

```console
docker exec -it cadence-server /bin/bash
```

```console
cadence --address $(hostname -i):7933 --do e-commerce domain register
```

4 - Start as many workers you'd like:

```
make cadence-worker
```

5 - Run the order workflow:

```
make order
```
