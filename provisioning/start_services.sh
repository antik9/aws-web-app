#!/bin/bash
nohup go run ../cmd/observer/main.go &
nohup go run ../cmd/notificator/main.go &
nohup go run ../cmd/server/main.go &
