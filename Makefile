# Harus sama dengan weaver.toml
GO_BUILD_OUT 	= unitedb
WEAVER_CONFIG 	= weaver.toml

.PHONY: generate

generate:
	weaver generate ./...

build: generate
	go build -o ${GO_BUILD_OUT} ./cmd/app

run-multi: build
	weaver multi deploy ${WEAVER_CONFIG}

run-single: build
	SERVICEWEAVER_CONFIG=${WEAVER_CONFIG} ./unitedb
