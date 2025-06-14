# Rock Paper Scissors Game Makefile

# Default values
PROJECT_ID ?= $(shell echo $$PROJECT_ID)
VERSION ?= latest
PLATFORM ?= linux/amd64

# Derived values
IMAGE_NAME = gcr.io/$(PROJECT_ID)/rock-paper-scissors-game:$(VERSION)

# Colors
BLUE = \033[0;34m
GREEN = \033[0;32m
YELLOW = \033[1;33m
NC = \033[0m

.PHONY: help build push deploy clean test run dev check-project health

help: ## Show this help message
	@echo "Rock Paper Scissors Game - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make deploy PROJECT_ID=my-project"
	@echo "  make build VERSION=v1.2.3 PROJECT_ID=my-project"
	@echo "  make auto-deploy PROJECT_ID=my-project"

check-project: ## Check if PROJECT_ID is set
	@if [ -z "$(PROJECT_ID)" ]; then \
		echo "$(YELLOW)Error: PROJECT_ID is required$(NC)"; \
		echo "Usage: make deploy PROJECT_ID=your-gcp-project-id"; \
		exit 1; \
	fi

build: check-project ## Build Docker image
	@echo "$(BLUE)Building image: $(IMAGE_NAME)$(NC)"
	docker buildx build --platform $(PLATFORM) -t $(IMAGE_NAME) .

push: check-project ## Build and push Docker image
	@echo "$(BLUE)Building and pushing image: $(IMAGE_NAME)$(NC)"
	docker buildx build --platform $(PLATFORM) -t $(IMAGE_NAME) . --push

deploy: push ## Build, push and deploy to Kubernetes
	@echo "$(BLUE)Deploying to Kubernetes...$(NC)"
	@sed -e 's/PROJECT_ID/$(PROJECT_ID)/g' -e 's/VERSION_TAG/$(VERSION)/g' deployment.yaml | kubectl apply -f -
	@echo "$(GREEN)Deployment complete!$(NC)"
	@kubectl get svc rock-paper-scissors-service
	@kubectl get pods -l app=rock-paper-scissors-game

auto-deploy: check-project ## Deploy with auto-generated version from git commit
	$(MAKE) deploy VERSION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "latest")

run: ## Run locally
	@echo "$(BLUE)Starting local server...$(NC)"
	go run .

dev: ## Run in development mode with auto-reload
	@echo "$(BLUE)Starting development server...$(NC)"
	@command -v air > /dev/null 2>&1 || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

test: ## Run tests
	go test ./...

clean: check-project ## Clean up Kubernetes resources
	@echo "$(YELLOW)Cleaning up resources...$(NC)"
	kubectl delete -f deployment.yaml --ignore-not-found=true

status: ## Show deployment status
	@echo "$(BLUE)Service status:$(NC)"
	@kubectl get svc rock-paper-scissors-service 2>/dev/null || echo "Service not found"
	@echo ""
	@echo "$(BLUE)Pod status:$(NC)"
	@kubectl get pods -l app=rock-paper-scissors-game 2>/dev/null || echo "No pods found"

logs: ## Show application logs
	@kubectl logs -l app=rock-paper-scissors-game --tail=50 -f

# Development helpers
format: ## Format Go code
	go fmt ./...

lint: ## Run linter
	@command -v golangci-lint > /dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

docker-run: build ## Run Docker container locally
	docker run --rm -p 8080:8080 $(IMAGE_NAME)

# Version management
release: ## Create a release (use VERSION=vX.Y.Z)
	@if [ "$(VERSION)" = "latest" ]; then \
		echo "$(YELLOW)Please specify VERSION for release: make release VERSION=v1.2.3$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Creating release $(VERSION)$(NC)"
	git tag $(VERSION)
	git push origin $(VERSION)
	$(MAKE) deploy VERSION=$(VERSION)

health: ## Run health check on deployed application
	@echo "$(BLUE)Running health check...$(NC)"
	@./health-check.sh
