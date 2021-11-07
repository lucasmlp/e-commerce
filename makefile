include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

deps:
	@ echo
	@ echo "Downloading dependencies..."
	@ echo
	@ go get -v ./...

update-deps:
	@ echo
	@ echo "Updating dependencies..."
	@ echo
	@ go get -u ./...

cadence-containers:
	@ echo
	@ echo "Starting Cassandra Cadence and Cadence Web..."
	@ echo
	@ docker-compose up -d cassandra cadence cadence-web

database:
	@ echo
	@ echo "Starting Cassandra Cadence and Cadence Web..."
	@ echo
	docker-compose -f stack.yml up

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
