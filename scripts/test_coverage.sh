#!/bin/sh

mkdir ./coverage 2>/dev/null

go test -coverprofile=coverage/cover.out ./...
go tool cover -html=coverage/cover.out

