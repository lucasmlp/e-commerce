include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

server:
	@ echo
	@ echo "Starting the API..."
	@ echo
	@ go run ./cmd/api/main.go

cadence-containers:
	@ echo
	@ echo "Starting Cassandra Cadence and Cadence Web..."
	@ echo
	@ docker-compose up -d cassandra cadence cadence-web

database:
	@ echo
	@ echo "Starting Mongo..."
	@ echo
	docker-compose up -d mongodb

cadence-worker:
	@ echo
	@ echo "Starting the Cadence Worker..."
	@ echo
	@ go run ./worker/main.go

order:
	@ go run ./trigger/main.go RunOrder

storage-check-reservation-finished:
	@ go run ./trigger/main.go StorageCheckReservationFinished $(filter-out $@,$(MAKECMDGOALS))

payment-finished:
	@ go run ./trigger/main.go PaymentFinished $(filter-out $@,$(MAKECMDGOALS))

shipment-finished:
	@ go run ./trigger/main.go ShipmentFinished $(filter-out $@,$(MAKECMDGOALS))
	
connect:
	@ kubectl proxy --port=9874 &