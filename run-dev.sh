#!/bin/bash
export GO_ENV=development
export LIQUIBASE_PROPERTIES=liquibase-dev.properties
export PORT=8081

# Kill any process using port 8081
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    lsof -ti :8081 | xargs kill -9 2>/dev/null || true
else
    # Linux
    fuser -k 8081/tcp 2>/dev/null || true
fi

go run github.com/cosmtrek/air@v1.49.0
