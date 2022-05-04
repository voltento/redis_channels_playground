start: build
	docker compose up -d

stop:
	docker compose down

start-db:
	docker compose up db

build: clean
	docker compose pull && docker compose up --build -d

start-with-logs: build
	docker compose up

clean: stop
	docker compose rm -f