include .env

BIN=out
DB_NAME=chirpy
DB_VOLUME=chirpy_data_vol
DB_IMAGE=postgres:17.4-alpine


.PHONY: run
run: build start-db
	./$(BIN)


.PHONY: build
build:
	go build -o $(BIN)


.PHONY: manage-db
manage-db: start-db
	@psql ${DB_URL}


.PHONY: migrate-up
migrate-up:
	@cd sql/schema && goose postgres ${DB_URL} up


.PHONY: migrate-down
migrate-down:
	@cd sql/schema && goose postgres ${DB_URL} down


.PHONY: start-db
start-db:
	@echo Starting postgres container
	-docker run \
		--detach \
		--name ${DB_NAME} \
		--env-file .env \
		-p 5432:5432 \
		-v ${DB_VOLUME}:/var/lib/postgresql/data \
		${DB_IMAGE}


.PHONY: stop-db
stop-db:
	@echo Stopping ${DB_NAME} container...
	-docker stop ${DB_NAME}


.PHONY: clean
clean:
	$(RM) $(BIN)


.PHONY: clean-db
clean-db: stop-db
	@echo Removing ${DB_NAME} container...
	-docker rm ${DB_NAME}


.PHONY: clean-db-vol
clean-db-vol: clean-db
	@echo Removing ${DB_NAME} volume...
	-docker volume rm ${DB_VOLUME}


.PHONY: clean-all
clean-all: clean clean-db clean-db-vol
	@echo Removing ${DB_IMAGE} image...
	-docker rmi ${DB_IMAGE}
