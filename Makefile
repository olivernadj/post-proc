.PHONY: start stop run createdb clean test

start:
	cd deployments && docker-compose up -d --build

stop:
	cd deployments && docker-compose stop

restart: | stop start

run:
	cd deployments && docker-compose up --build

createdb: | restart
	docker exec -t post_proc_db sh -c /scripts/create.sql.sh

clean: | stop
	cd deployments && docker-compose rm -f

test:
	cd test/deployments &&\
	docker-compose stop &&\
	docker-compose rm -f &&\
	docker-compose up -d --build db_test &&\
	docker exec -t post_proc_db_test sh -c /scripts/create.sql.sh &&\
	docker-compose up --build restapi_test &&\
	docker-compose up --build statemachine_test &&\
	docker-compose stop &&\
	docker-compose rm -f

