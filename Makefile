.PHONY: build run test clean docker-build docker-run k8s-deploy run-examples

# Build all binaries
build:
	go build -o bin/api-server ./cmd/api-server
	go build -o bin/chaos-operator ./cmd/chaos-operator
	go build -o bin/chaos-cli ./cmd/cli

# Run the API server locally
run-api:
	go run ./cmd/api-server/main.go

# Run the chaos operator locally
run-operator:
	go run ./cmd/chaos-operator/main.go

# Run the entire platform locally
run:
	@echo "Starting API server..."
	go run ./cmd/api-server/main.go & \
	@echo "Starting chaos operator..."
	go run ./cmd/chaos-operator/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker images
docker-build:
	docker build -t chaos-platform/api-server:latest -f deployments/docker/api-server.Dockerfile .
	docker build -t chaos-platform/chaos-operator:latest -f deployments/docker/chaos-operator.Dockerfile .

# Run with Docker Compose
docker-run:
	docker-compose up

# Deploy to Kubernetes
k8s-deploy:
	kubectl apply -f deployments/kubernetes/secrets.yaml
	kubectl apply -f deployments/kubernetes/postgres.yaml
	kubectl apply -f deployments/kubernetes/api-server.yaml
	kubectl apply -f deployments/kubernetes/chaos-operator.yaml
	kubectl apply -f deployments/kubernetes/monitoring.yaml
	kubectl apply -f deployments/kubernetes/ingress.yaml

# Initialize the database schema
init-db:
	go run scripts/init_db.go

# Run examples
run-examples:
	bash examples/external-target/test-external-target.sh