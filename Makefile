ENV_PATH = $(PWD)/.env
include $(ENV_PATH)

DOWN ?= def
UP ?= def
FORCE ?=def

run-dev:
	go run cmd/main.go
run-prod:
	./bin/macos-amd64
build:
	GOOS=darwin GOARCH=amd64 cd cmd && go build -o ../bin/macos-amd64
build-win:
	GOOS=windows GOARCH=amd64 cd cmd && go build -o ../bin/windows-amd64.exe
swag:
	swag init -g cmd/main.go
migrate-up:
	migrate -path ./migrations -database '$(DB_URL)?sslmode=disable' up $(UP)
migrate-down:
	migrate -path ./migrations -database '$(DB_URL)?sslmode=disable' down $(DOWN)
migrate-force:
	migrate -path ./migrations -database '$(DB_URL)?sslmode=disable' force $(FORCE)


#  migrate -path ./migrations -database 'postgres://admin:maxim@localhost:5432/go-fiber?sslmode=disable' up 1
#  migrate -path ./migrations -database 'postgres://admin:maxim@localhost:5432/go-fiber?sslmode=disable' down 1


# add -e ENV="dev" or "prod" for migrations
# migrate create -ext sql -dir ./migrations -seq (name)
