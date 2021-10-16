#!/bin/bash
go test ./... -coverprofile=coverage.out -bench . -count=1
go tool cover -html=coverage.out -o coverage.html