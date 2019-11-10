.PHONY: all

all: server

server:
	go build -o web-app cmd/server/main.go
