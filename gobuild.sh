#!/bin/bash
for r in route/*; do
    if [ -d "$r" ]; then
        r=$(basename "$r")
        env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$r route/$r/main.go
    fi
done
