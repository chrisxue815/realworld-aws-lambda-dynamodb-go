#!/bin/bash
for r in route/*; do
    if [ -d "$r" ]; then
        r=$(basename "$r")
        env GOOS=linux go build -ldflags="-s -w" -o bin/$r route/$r/main.go
    fi
done
