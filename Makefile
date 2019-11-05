.PHONY: start stop run createdb clean

start:
	cd deployments && docker-compose up -d --build

stop:
	cd deployments && docker-compose stop

restart: | stop start

run:
	cd deployments && docker-compose up --build

createdb:
	docker exec -t post_proc_db sh -c /scripts/create.sql.sh

clean:
	cd deployments && docker-compose rm -f