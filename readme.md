# School Payment Backend Service

## Overview

A Go-based backend service for managing school payments, built with Fiber framework and modern Go practices. This service handles student payments, billing management, and various school-related financial operations with support for multiple payment methods and comprehensive audit trails.

## Core Technologies

- Go (Golang) 1.21.5
- Fiber (Web Framework)
- GORM (ORM)
- PostgreSQL (Database)
- Liquibase (Database Migration)
- JWT (Authentication)
- Docker
- Midtrans (Payment Gateway)

## Prerequisites

- Go (Golang) version 1.21.5 or higher
- Git
- PostgreSQL 12+
- Java Runtime Environment (JRE) 11+ (for Liquibase)
- Liquibase 4.23.0
- Docker (optional, for containerization)

### Installing Liquibase

1. **Windows**:

   - Download Liquibase 4.23.0 from [official website](https://www.liquibase.org/download)
   - Add Liquibase to your system PATH

2. **Linux/MacOS**:
   ```bash
   wget https://github.com/liquibase/liquibase/releases/download/v4.23.0/liquibase-4.23.0.zip
   unzip liquibase-4.23.0.zip -d /opt/liquibase
   export PATH=$PATH:/opt/liquibase
   ```

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://gitlab.com/gudangsolusi/school-payment-backend
   cd school-payment-backend
   ```

2. **Set up environment variables:**

   ```bash
   # For production
   cp .env.example .env
   # Edit .env with your production configuration

   # For development
   cp .env.dev.example .env.dev
   # Edit .env.dev with your development configuration
   ```

3. **Install dependencies:**

   ```bash
   go mod download
   ```

4. **Run database migrations:**
   ```bash
   liquibase update
   ```

## Running the Application

### Development Mode

For Linux/MacOS:

```bash
./run-dev.sh
```

For Windows:

```bash
./run-dev.bat
```

Or manually:

```bash
# For Linux/MacOS
export GO_ENV=development
go run main.go

# For Windows
set GO_ENV=development
go run main.go
```

### Production Mode

```bash
go build -o app
./app
```

### Docker Deployment

```bash
docker build -t school-payment-backend .
docker run -p 8081:8081 school-payment-backend
```

## Project Structure

```
├── configs/         # Configuration and database setup
├── controllers/     # HTTP request handlers
├── models/         # Data models
├── repositories/   # Database operations
├── services/       # Business logic
├── utilities/      # Helper functions
├── routes/         # API route definitions
├── dtos/          # Data Transfer Objects
└── db/
    └── changelog/  # Liquibase migration files
```

## Key Features

- User Management (Role-based authentication)
- School Management
- Student Management
- Payment Processing
  - Multiple payment methods support
  - Midtrans payment gateway integration
  - Transaction history tracking
  - Payment status monitoring
- Billing Management
  - Multiple billing types
  - Flexible payment periods
  - Bank account management
- Transaction History
- Report Generation
- Email Notifications
- Audit Trail System
- File Upload Handling

## Database Migration

The project uses Liquibase for database version control. Migration files are located in `db/changelog/`.

To run migrations:

```bash
liquibase update
```

To rollback last migration:

```bash
liquibase rollbackCount 1
```

## Testing

The project provides several options for running tests:

### Running All Tests

Using shell script (Linux/MacOS):

```bash
./run-test.sh
```

Using batch script (Windows):

```bash
./run-test.bat
```

### Running Specific Package Tests

Using shell script (Linux/MacOS):

```bash
./run-test-specific.sh
```

Using batch script (Windows):

```bash
./run-test-specific.bat
```

This will present a menu to select which packages to test:

1. controllers
2. services
3. repositories
4. routes
5. utilities
6. models

All test methods will:

- Clean previous coverage files
- Run tests with coverage for selected packages
- Display coverage report in terminal
- Generate HTML coverage report (coverage.html)

### Running Standard Go Tests

You can also run tests using standard Go commands:

```bash
go test ./test/
```

## Contributing

Please follow the merge request template located in `.gitlab/merge_request_templates/aigs_mr_template.md` when submitting changes.

```

```
