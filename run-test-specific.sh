#!/bin/bash

# Function to display available packages
show_packages() {
    echo "Available packages:"
    echo "1. controllers"
    echo "2. services"
    echo "3. repositories"
    echo "4. routes"
    echo "5. utilities"
    echo "6. models"
    echo "0. Exit"
}

# Show packages and get user input
show_packages
echo "Enter package number (or multiple numbers separated by space, e.g., '1 2 3'):"
read -r selection

# Exit if user selected 0
if [[ $selection == "0" ]]; then
    exit 0
fi

# Convert selection to package paths
packages=""
for num in $selection; do
    case $num in
        1) packages="$packages ./controllers/...";;
        2) packages="$packages ./services/...";;
        3) packages="$packages ./repositories/...";;
        4) packages="$packages ./routes/...";;
        5) packages="$packages ./utilities/...";;
        6) packages="$packages ./models/...";;
        *) echo "Invalid selection: $num";;
    esac
done

if [ -z "$packages" ]; then
    echo "No valid packages selected."
    exit 1
fi

# Clean previous coverage files
rm -f coverage.out

# Run tests with coverage for selected packages
echo "Running tests for selected packages: $packages"
go test $packages -coverprofile=coverage.out

# Display coverage report in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Print completion message
echo "Testing completed. Coverage report generated in coverage.html"