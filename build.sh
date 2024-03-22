#!/bin/bash

echo "Download dependancies"
go mod tidy
echo "build the binary"
GOARCH=arm64 GOOS=linux go build -o bootstrap main.go
echo "Create a zip file"
zip awslambda.zip bootstrap
echo "Cleaning up"
rm -r bootstrap