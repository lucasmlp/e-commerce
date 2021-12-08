include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

server:
	@ echo
	@ echo "Starting the API..."
	@ echo
	@ go run ./cmd/api/main.go

test:
	@ echo
	@ echo "Starting running tests..."
	@ echo
	@ go clean -testcache & go test -cover ./...

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
	@ go run ./cmd/workflowworker/main.go

order:
	@ go run ./cadence/trigger/main.go RunOrder

storage-check-reservation-finished:
	@ go run ./trigger/main.go StorageCheckReservationFinished $(filter-out $@,$(MAKECMDGOALS))

payment-finished:
	@ go run ./trigger/main.go PaymentFinished $(filter-out $@,$(MAKECMDGOALS))

shipment-finished:
	@ go run ./trigger/main.go ShipmentFinished $(filter-out $@,$(MAKECMDGOALS))
	
connect:
	@ kubectl proxy --port=9874 &

mock:
	@ echo
	@ echo "Starting building mocks..."
	@ echo
	@ mockgen -source=domain/orders/service.go -destination=domain/orders/mocks/service_mock.go -package=mocks
	@ mockgen -source=domain/orders/repository.go -destination=domain/orders/mocks/repository_mock.go -package=mocks
	@ mockgen -source=domain/products/service.go -destination=domain/products/mocks/service_mock.go -package=mocks
	@ mockgen -source=domain/products/repository.go -destination=domain/products/mocks/repository_mock.go -package=mocks