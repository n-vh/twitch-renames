APPS=api irc server
LOCAL_PATH=build/local
PROD_PATH=build/prod

dev-api:
	go run ./cmd/api/main.go

dev-irc:
	go run ./cmd/irc/main.go

dev-server:
	go run ./cmd/server/main.go

local:
	rm -rf -f ./${LOCAL_PATH}/*
	for i in ${APPS}; do \
		go build -o ./${LOCAL_PATH}/$$i.exe ./cmd/$$i/main.go; \
	done

prod:
	rm -rf -f ./${PROD_PATH}/*
	for i in ${APPS}; do \
		env GOOS=linux GOARCH=amd64 go build -o ./${PROD_PATH}/$$i ./cmd/$$i/main.go; \
	done