SERVER_BINARY=sport4all
DOCKER_DIR=docker

DOCUMENTATION_CONTAINER_NAME=documentation

# code
run-linter:
	golangci-lint -c .golangci.yml run ./...

format:
	go fmt github.com/SIBIRSKAYA-KORONA/sport4all-backend/...

generate:
	go generate ./...


# application
build-server:
	go build -o ${SERVER_BINARY} ./cmd

run:
	go run ./cmd -c conf/config.yml

# docs
generate-swagger:
	swagger generate spec -o ./docs/swagger.yaml --scan-models

serve-swagger:
	swagger serve -F=swagger ./docs/swagger.yaml

doc-host:
	docker run --name=documentation -d -p 5757:8080 -e SWAGGER_JSON=/swagger.yaml -v ${CURDIR}/docs/swagger.yaml:/swagger.yaml swaggerapi/swagger-ui

doc-stop:
	docker stop ${DOCUMENTATION_CONTAINER_NAME}
	docker rm ${DOCUMENTATION_CONTAINER_NAME}

# docker
docker-make-all-images: docker-make-builder-image docker-make-api-image

docker-make-builder-image:
	docker build -t sport-builder -f ${DOCKER_DIR}/builder.Dockerfile .

docker-make-api-image:
	docker build -t sport-api -f ${DOCKER_DIR}/Dockerfile .

# docker-compose
start:
	docker-compose -f ${DOCKER_DIR}/docker-compose.yml up -d

stop:
	docker-compose -f ${DOCKER_DIR}/docker-compose.yml stop

down:
	docker-compose -f ${DOCKER_DIR}/docker-compose.yml down