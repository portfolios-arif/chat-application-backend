#!/bin/bash

# Get project root dir
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
go run "$PROJECT_ROOT/cmd/main.go"
