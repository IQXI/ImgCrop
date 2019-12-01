UTS=internal\lru_cache
ITS=internal\integration_tests
ICD=cmd
DD=docker

all: build test run
test: unit_test integration_test

unit_test:
@cd $(UTS)
@echo "UTC(lru_cache)"
go test -v
integration_test:
@cd $(ITS)
@echo "ITS"
docker-compose -f $(DD)\docker-compose.test.yaml up tests
docker-compose -f $(DD)\docker-compose.test.yaml down

build:
@cd $(ICD)
@echo "Build ImgCrop"
@echo "Go vet"
@go vet
@echo "Go lint"
golangci-lint run
go build -race -v cmd/main.go
run:
docker-compose -f $(DD)\docker-compose.yaml up imgcrop