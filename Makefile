.PHONY: up down build logs test clean

# Starts the application in detached mode, building images if necessary
up:
	@echo "Starting Kite infrastructure..."
	docker compose up --build -d
	@echo "Frontend is running at http://localhost:5173"

# Stops the application and removes the containers
down:
	@echo "Stopping Kite infrastructure..."
	docker compose down

# Shuts down the application and completely wipes the database volume
clean:
	@echo "Wiping database and stopping containers..."
	docker compose down -v

# Tails the logs of all containers
logs:
	docker compose logs -f

# Runs the backend test suite locally
test:
	@echo "Running backend unit tests..."
	cd backend && go test -v ./...