package configs

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func RunLiquibaseUpdate() error {

	liquibaseProperties := os.Getenv("LIQUIBASE_PROPERTIES")
	if liquibaseProperties == "" {
		log.Fatal("LIQUIBASE_PROPERTIES not set in environment variables")
	}

	cmd := exec.Command("liquibase", "update", "--defaultsFile="+liquibaseProperties)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error during migration:", string(output))
		fmt.Println("Attempting rollback...")
		rollbackErr := RollbackLastMigration()
		if rollbackErr != nil {
			return fmt.Errorf("rollback failed: %v", rollbackErr)
		}
		return fmt.Errorf("migration failed: %v", err)
	}

	fmt.Println("Liquibase migration successful")
	return nil
}

func RollbackLastMigration() error {
	liquibaseProperties := os.Getenv("LIQUIBASE_PROPERTIES")
	if liquibaseProperties == "" {
		log.Fatal("LIQUIBASE_PROPERTIES not set in environment variables")
	}

	cmd := exec.Command("liquibase", "rollbackCount", "1", "--defaultsFile="+liquibaseProperties)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error during rollback:", string(output))
		return err
	}

	fmt.Println("Rollback successful")
	return nil
}
