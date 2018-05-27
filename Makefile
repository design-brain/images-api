GO_FILES := ./...
SERVICE_NAME := images-api
export

.PHONY: all compose-deps compose-down compose-integration-test compose-up docker-lint gen install lint run run-dotenv test vendor

all: install

compose-deps:
	docker-compose up -d cache
	docker-compose up -d db

compose-down:
	docker-compose down

compose-integration-test: compose-deps
	docker-compose build api
	docker-compose run --entrypoint "make integration-test" api

compose-up: compose-deps
	docker-compose up -d --build api

docker-lint:
	docker build -f Dockerfile.lint -t $(SERVICE_NAME)-lint .
	docker run $(SERVICE_NAME)-lint

gen:
	retool do protoc --proto_path=. --twirp_out=$(GOPATH)/src --go_out=$(GOPATH)/src ./rpc/images/images.proto

install:
	go install -v ./...

lint:
	gometalinter --config=.gometalinter.json --deadline=2m --exclude=rpc --exclude=vendor $(GO_FILES)

run: install
	$(SERVICE_NAME)

run-dotenv: install
	$(SERVICE_NAME) -dotenv

test:
	go test -v $(GO_FILES)

vendor:
	retool do vgo build ./...
	retool do vgo vendor
