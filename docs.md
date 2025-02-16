# School Payment Backend Service

## Project Overview

This is a Go-based backend service for a school payment management system. The project uses modern Go practices and several key technologies:

### Core Technologies

- Go (Golang) 1.21.5
- Fiber (Web Framework)
- GORM (ORM)
- PostgreSQL (Database)
- JWT (Authentication)
- Liquibase (Database Migration)
- Docker

### Key Features

1. **User Management**

   - Role-based authentication
   - Email verification
   - Password management
   - User blocking capabilities

2. **School Management**

   - School registration and profile management
   - Class and grade management
   - Student management
   - Parent/Guardian management

3. **Payment System**

   - Multiple payment methods support
   - Transaction processing
   - Billing management
   - Payment status tracking
   - Invoice generation (PDF)

4. **Additional Features**
   - Email notifications
   - File upload handling
   - Audit trail logging
   - Webhook handling for payment updates

## Architecture

### Project Structure

The project follows a clean architecture pattern with clear separation of concerns:

```
├── configs/         # Configuration and database setup
├── controllers/     # HTTP request handlers
├── models/         # Data models
├── repositories/   # Database operations
├── services/       # Business logic
├── utilities/      # Helper functions
├── routes/         # API route definitions
└── dtos/          # Data Transfer Objects
```

### Key Components

1. **Authentication System**

```26:76:services/loginService.go
func (loginService *LoginService) LoginService(email string, password string, firebaseToken string) (string, error) {
	var user models.User
	var err error
	platform := "website"
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailPattern)
	if re.MatchString(email) {
		user, err = repositories.GetUserByEmail(email)
		if err != nil {
			return "", fmt.Errorf("Email/Username atau password yang Anda masukkan salah. Silahkan coba lagi")
		}
	} else {

		user, err = repositories.GetUserByUsername(email)
		if err != nil {
			return "", fmt.Errorf("Email/Username atau password yang Anda masukkan salah. Silahkan coba lagi")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("Email/Username atau password yang Anda masukkan salah. Silahkan coba lagi")
	}

	// Step 3: Generate JWT token
	token, err := utilities.GenerateJWT(user, firebaseToken)
	if err != nil {
		return "", fmt.Errorf("Oops something wrong.")
	}

	if firebaseToken != "" {
		platform = "mobile"
	}

	//save to audit trail
	newAuditTrail := models.AuditTrail{
		UserID:     user.ID,
		Email:      user.Email,
		Role:       user.Role.Name,
		UserAction: "Login",
		ApiPath:    "/login",
		LogTime:    time.Now(),
		Platform:   platform,
		FirebaseID: firebaseToken,
	}

	newAuditTrail.CreatedBy = int(user.ID)
	newAuditTrail.UpdatedBy = int(user.ID)
	err = loginService.auditTrailRepository.CreateDataAuditTrail(&newAuditTrail)
	return token, nil

```

The authentication system uses JWT tokens and includes email verification.

2. **Transaction Processing**

```239:276:services/transactionService.go

	if dataTransaction != nil {
		history := models.TransactionBillingHistory{
			TransactionBillingId: dataTransaction.ID,
			// TransactionDate:      time.Now(),
			// TransactionAmount:    request.TotalAmount,
			ReferenceNumber: referenceNumber,
			// TransactionType:      "kasir",
			// Description:          request.Description,
			OrderID:           "",
			InvoiceNumber:     invoiceNumber,
			TransactionStatus: "PS02",
		}
		history.CreatedBy = userId
		_, err = transactionService.transactionRepository.CreateTransactionHistoryRepository(tx, &history)
		if err != nil {
			// tx.Rollback()
			return models.TransactionBilling{}, err
		}
	}

	billingStudentIds := request.BillingStudentIds
	var billingStudentIdsInt []int
	for _, idStr := range billingStudentIds {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			// Tangani error jika konversi gagal
			fmt.Println("Error konversi:", err)
			return models.TransactionBilling{}, err
		}
		billingStudentIdsInt = append(billingStudentIdsInt, idInt)
	}
	if dataTransaction != nil {
		err := repositories.UpdateStatusPayment(billingStudentIdsInt)
		if err != nil {
			return models.TransactionBilling{}, err
		}
	}
```

Handles payment processing with status tracking and history logging.

3. **Email Service**

```111:128:utilities/verificationEmail.go

	// Create a new email message
	m := mail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) //to send plain text use text/plain

	// Create a new SMTP dialer
	d := mail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
```

Handles email notifications using SMTP with HTML templates.

## Security Features

1. **Role-Based Access Control**

```252:266:controllers/userController.go
func (userController *UserController) CheckAccessToUser(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	_, err := services.GetRoleByID(uint(roleID))
	if err != nil {
		return err
	}

	if roleID == 1 || roleID == 5 {
		return nil
	}

	return fmt.Errorf("User can't access this page")
}
```

2. **Password Security**

- Bcrypt hashing for passwords
- Password validation rules
- Secure password reset flow

## Deployment

The project includes Docker configuration for containerization:

```1:38:Dockerfile
# FROM registry.gitlab.com/gudangsolusi/digiform-images/golang:alpine3.13 AS build-env
FROM golang:1.21.5-alpine

ENV GOPROXY=direct
ENV GO111MODULE "auto"
ENV PATH $PATH:$HOME/go/bin

# Install necessary packages (git, Java, wget, unzip, tzdata)
RUN apk update && apk add --no-cache bash git openjdk11 wget unzip tzdata

# Set timezone to Asia/Jakarta
ENV TZ="Asia/Jakarta"
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Optional: Set JAVA_HOME environment variable
ENV JAVA_HOME /usr/lib/jvm/java-11-openjdk
ENV PATH $JAVA_HOME/bin:$PATH

# Install Liquibase
RUN wget https://github.com/liquibase/liquibase/releases/download/v4.23.0/liquibase-4.23.0.zip
RUN unzip liquibase-4.23.0.zip -d /opt/liquibase
ENV PATH $PATH:/opt/liquibase

# Application setup
ADD . /main/
WORKDIR /main
# COPY .env.temp .env

# Install swag CLI and generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Fetch Go dependencies and build
RUN go get -d ./...
RUN go build -o app .

# Increase file descriptors limit
RUN ulimit -n 65536
```

Key deployment features:

- Alpine-based lightweight container
- Timezone configuration (Asia/Jakarta)
- Liquibase integration for database migrations
- Swagger documentation generation
- Automatic environment setup

Would you like me to elaborate on any specific aspect of the project?
 