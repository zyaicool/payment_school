@echo off
echo Running tests...

:: Clean previous coverage files
del /f coverage.out coverage.html 2>nul

:: Run tests with coverage for specific packages
go test ./controllers/... ./services/... ./repositories/... ./routes/... ./utilities/... ./models/... -coverprofile=coverage.out

:: Display coverage report in terminal
go tool cover -func=coverage.out

:: Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

:: Print completion message
echo Testing completed. Coverage report generated in coverage.html