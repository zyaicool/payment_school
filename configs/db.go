package configs

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to database")
	}
}

func SynchDB() {
	fmt.Println("Database schema managed by Liquibase.")
}

func ExecuteSQLScriptUsingGORM(db *gorm.DB, filePath string) error {
	// Open the SQL script file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file line by line and build the SQL query
	scanner := bufio.NewScanner(file)
	var queryBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments and empty lines
		if strings.HasPrefix(strings.TrimSpace(line), "--") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		queryBuilder.WriteString(line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Execute the SQL query using GORM's Exec method
	query := queryBuilder.String()
	if err := db.Exec(query).Error; err != nil {
		return fmt.Errorf("failed to execute SQL script: %v", err)
	}

	fmt.Println("SQL script executed successfully")
	return nil
}
