run:
	go run cmd/main.go

up:
	docker compose up -d

build:
	docker compose build

up_build:
	docker compose up -d --build

stop:
	docker compose stop

down:
	docker compose down

clean:
	docker volume rm auth-service_postgres_data