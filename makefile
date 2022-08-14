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
	
kubectl-proxy:
	@ kubectl proxy --port=9874

mock:
	@ echo
	@ echo "Starting building mocks..."
	@ echo
	@ mockgen -source=domain/orders/service.go -destination=domain/orders/mocks/service_mock.go -package=mocks
	@ mockgen -source=domain/orders/repository.go -destination=domain/orders/mocks/repository_mock.go -package=mocks
	@ mockgen -source=domain/products/service.go -destination=domain/products/mocks/service_mock.go -package=mocks
	@ mockgen -source=domain/products/repository.go -destination=domain/products/mocks/repository_mock.go -package=mocks

docker-image:
	@ echo
	@ echo "Building docker image..."
	@ echo
	@ docker build -t machado-br/e-commerce-api:latest .

docker-tag-aws:
	@ echo
	@ echo "Tagging docker image for AWS..."
	@ echo
	@ docker tag machado-br/e-commerce-api:latest 774429751797.dkr.ecr.us-west-2.amazonaws.com/e-commerce-api:latest

login-aws-ecr:
	@ echo
	@ echo "Logging in AWS ECR..."
	@ echo
	@ aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 774429751797.dkr.ecr.us-west-2.amazonaws.com

docker-push-aws:
	@ echo
	@ echo "Pushing docker image to AWS ECR..."
	@ echo
	@ docker push 774429751797.dkr.ecr.us-west-2.amazonaws.com/e-commerce-api:latest