@echo off
setlocal enabledelayedexpansion

:show_packages
echo Available packages:
echo 1. controllers
echo 2. services
echo 3. repositories
echo 4. routes
echo 5. utilities
echo 6. models
echo 0. Exit

set /p selection="Enter package number (or multiple numbers separated by space, e.g., '1 2 3'): "

if "%selection%"=="0" exit /b

:: Initialize empty packages string
set "packages="

:: Process each selected number
for %%n in (%selection%) do (
    if "%%n"=="1" set "packages=!packages! ./controllers/..."
    if "%%n"=="2" set "packages=!packages! ./services/..."
    if "%%n"=="3" set "packages=!packages! ./repositories/..."
    if "%%n"=="4" set "packages=!packages! ./routes/..."
    if "%%n"=="5" set "packages=!packages! ./utilities/..."
    if "%%n"=="6" set "packages=!packages! ./models/..."
)

if "%packages%"=="" (
    echo No valid packages selected.
    exit /b 1
)

:: Clean previous coverage files
del /f coverage.out coverage.html 2>nul

:: Run tests with coverage for selected packages
echo Running tests for selected packages: %packages%
go test %packages% -coverprofile=coverage.out

:: Display coverage report in terminal
go tool cover -func=coverage.out

:: Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

:: Print completion message
echo Testing completed. Coverage report generated in coverage.html
