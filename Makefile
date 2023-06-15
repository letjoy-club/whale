build:
	go build -o bin/server cmd/server/main.go

linux-server:
	GOOS=linux GOARCH=amd64 go build -o bin/whale cmd/server/main.go

gen:
	go run cmd/gqlgen/main.go generate

db-init:
	go run cmd/db/main.go -init

db-init-prod:
	go run cmd/db/main.go -init -conf conf.prod.yaml

db-gen:
	go run cmd/db/main.go

tools:
	go build -o bin/auth-tool cmd/auth-tool/main.go

smew-genqlient:
	cd pkg/gqlient/smew && make schema && make gen

hoopoe-genqlient:
	cd pkg/gqlient/hoopoe && make schema && make gen

scream-genqlient:
	cd pkg/gqlient/scream && make schema && make gen

.PHONY: build
