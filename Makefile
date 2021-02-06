run-linter:
	golangci-lint -c .golangci.yml run ./...

format:
	go fmt github.com/SIBIRSKAYA-KORONA/sport4all-backend/...

run:
	go run cmd/main.go -c cmd/config.yml

generate:
	go generate ./...

generate-swagger:
	swagger generate spec -o ./docs/swagger.yaml --scan-models

serve-swagger:
	swagger serve -F=swagger ./docs/swagger.yaml