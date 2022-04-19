APP_NAME=heyui
VERSION=1.0
KEYFOLDER:=key-pair
CURFOLDER:=$(shell pwd)

GITCOMMITCOUNT:=$$(git rev-list HEAD | wc -l | tr -d ' ')
GITHASH:=$$(git rev-parse --short HEAD)
DATETIME:=$$(date "+%Y%m%d%H%M%S")
VERSIONS:=$(VERSION).$(GITCOMMITCOUNT)-$(DATETIME)-$(GITHASH)

.PHONY: keygen postgres
default: build

keygen: clean
	mkdir $(KEYFOLDER)
	openssl genrsa -out $(KEYFOLDER)/id_rsa 4096
	openssl rsa -in $(KEYFOLDER)/id_rsa -pubout -out $(KEYFOLDER)/id_rsa.pub

postgres:
#	docker stop postgres
#	docker rm postgres
	chmod +x docker-entrypoint-initdb.d/init.sh
	docker run -d \
		--name postgres \
		-p 5432:5432 \
		-e POSTGRES_USER=root \
		-e POSTGRES_PASSWORD=password \
		-e APP_DB_USER=ui_test \
		-e APP_DB_PASS=ui_test \
		-e APP_DB_NAME=ui_test \
		-e PGDATA=/var/lib/postgresql/data/pgdata \
		-v $(CURFOLDER)/docker-entrypoint-initdb.d:/docker-entrypoint-initdb.d \
		postgres:14.2

clean:
	go clean
	rm -rvf $(KEYFOLDER)

doc:
	swag init --parseInternal --parseDependency

build: clean doc keygen
	go mod download && go mod verify
	go test -race ./...
	go build -v -o server -gcflags='-N -l' -ldflags "-X main.ServiceVersion=$(VERSION)" main.go
