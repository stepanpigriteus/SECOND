PROJECT_NAME=1337b04rd

up:
	docker-compose up -d --build

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f

ps:
	docker-compose ps

init:
	chmod +x scripts/init-bucket.sh

rebuild:
	docker-compose up -d --build --force-recreate

clean:
	docker-compose down -v --remove-orphans

status: ps

.PHONY: up down restart logs ps init rebuild clean status
