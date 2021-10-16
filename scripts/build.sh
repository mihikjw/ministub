#!/bin/bash
APP_NAME=ministub
UNAME=$(uname -s | awk '{print tolower($0)}')

CGO_ENABLED=0 GOOS=${UNAME} go build -o bin/${APP_NAME} -v cmd/ministub.go