PG_USERNAME = meli_user
PG_PASSWORD = qJDxWtxBvD3LR
PG_DB_NAME = meli
PG_PORT = 5432
PG_CONTAINER_NAME = meli-postgres
IP_POSTGRES = $(shell docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' meli-postgres)

build:
	./launcher.sh

format:
	echo "Formatting all code..."
	@docker run -v golang -v ${PWD}:/app -v ${PWD}/tmp:/go/pkg/mod -w /app golang gofmt -s -w .

run:
	@DB_HOST=${IP_POSTGRES} DB_USER=${PG_USERNAME} DB_PASS=${PG_PASSWORD} DB_NAME=${PG_DB_NAME} DB_PORT=5432 \
		CREDENTIALS_JSON_GMAIL=$(PWD)/credentials.json ./appBuilt

postgres-run:
	@echo "Running PostgreSQL ..."
	@docker run --name ${PG_CONTAINER_NAME} \
		-v ${PWD}/postgres:/docker-entrypoint-initdb.d \
		-e POSTGRES_PASSWORD=${PG_PASSWORD} \
		-e POSTGRES_USER=${PG_USERNAME} \
		-e POSTGRES_DB=${PG_DB_NAME} \
		-d postgres

postgres-stop:
	@echo "Stoping PostgreSQL Database..."
	@docker stop meli-postgres
	@docker rm meli-postgres

postgres-logs:
	@docker logs -t meli-postgres

postgres-data:
	@docker exec -i ${PG_CONTAINER_NAME} bash -c "PGPASSWORD=${PG_PASSWORD} psql -U ${PG_USERNAME} -d ${PG_DB_NAME} -c 'SELECT * FROM challenge'"