@echo off
set GO_ENV=development
set LIQUIBASE_PROPERTIES=liquibase-dev.properties
set PORT=8081
netstat -ano | findstr :8081 > nul && (
    FOR /F "tokens=5" %%P IN ('netstat -ano ^| findstr :8081') DO taskkill /F /PID %%P
)
go run github.com/cosmtrek/air@v1.49.0
