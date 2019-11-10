.PHONY: all

all: server observer notificator

server:
	go build -o web-app cmd/server/main.go

observer:
	go build -o ip-observer cmd/observer/main.go

notificator:
	go build -o file-notificator cmd/notificator/main.go
