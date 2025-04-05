package migrations

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// tableColumnMap defines which columns should be made nullable for each table
var tableColumnMap = map[string][]string{
	"leaderboards": {
		"description",
	},
	"participants": {
		"external_id",
		"metadata",
	},
	"leaderboard_entries": {
		// All non-null constraints are actually required
	},
	"leaderboard_metrics": {
		// All non-null constraints are actually required
	},
}

// requiredFieldsWithDefaults defines required fields that need default values for existing NULL data
var requiredFieldsWithDefaults = map[string]map[string]string{
	"leaderboards": {
		"name":             "Unnamed Leaderboard",
		"category":         "general",
		"type":             "individual",
		"time_frame":       "all-time",
		"sort_order":       "descending",
		"visibility_scope": "public",
		"is_active":        "true",
	},
	"participants": {
		"name": "Unnamed Participant",
		"type": "individual",
	},
	"leaderboard_entries": {
		"rank":         "0",
		"score":        "0.0",
		"last_updated": "NOW()",
	},
	"leaderboard_metrics": {
		"weight":           "1.0",
		"display_priority": "0",
	},
}

// fieldTypes maps column names to their expected types for creating missing columns
var fieldTypes = map[string]string{
	"name":             "VARCHAR(255)",
	"description":      "TEXT",
	"category":         "VARCHAR(255)",
	"type":             "VARCHAR(255)",
	"time_frame":       "VARCHAR(255)",
	"sort_order":       "VARCHAR(255)",
	"visibility_scope": "VARCHAR(255)",
	"is_active":        "BOOLEAN",
	"external_id":      "VARCHAR(255)",
	"metadata":         "JSONB",
	"rank":             "INTEGER",
	"score":            "FLOAT8",
	"last_updated":     "TIMESTAMP",
	"weight":           "FLOAT8",
	"display_priority": "INTEGER",
}

// Migration01FixSchema handles various schema fixes including:
// 1. Fix CreatedAt and UpdatedAt field tags
// 2. Make specified fields nullable across all tables
// 3. Fix required fields with default values
func Migration01FixSchema(db *gorm.DB) error {
	fmt.Println("Running Migration01FixSchema...")

	// Fix base model timestamp columns in all tables
	if err := fixBaseModelColumns(db); err != nil {
		return err
	}

	// Process each table's specified nullable columns
	for table, columns := range tableColumnMap {
		if len(columns) == 0 {
			continue
		}

		fmt.Printf("Processing table '%s' for nullable columns...\n", table)

		// Check if table exists
		if !tableExists(db, table) {
			fmt.Printf("Table '%s' doesn't exist yet, skipping\n", table)
			continue
		}

		for _, column := range columns {
			if err := makeColumnNullable(db, table, column); err != nil {
				return err
			}
		}
	}

	// Fix required fields that have NULL values or create if they don't exist
	for table, fields := range requiredFieldsWithDefaults {
		// First ensure the table exists
		if !tableExists(db, table) {
			fmt.Printf("Table '%s' doesn't exist yet, creating...\n", table)
			// If it's a core table but doesn't exist yet, create a minimal version
			if err := createTableIfNeeded(db, table); err != nil {
				return err
			}
		}

		fmt.Printf("Processing table '%s' for required fields...\n", table)

		for field, defaultVal := range fields {
			if err := ensureColumnWithDefault(db, table, field, defaultVal); err != nil {
				return err
			}
		}
	}

	fmt.Println("Migration01FixSchema completed successfully")
	return nil
}

// createTableIfNeeded creates a minimal version of core tables if they don't exist
func createTableIfNeeded(db *gorm.DB, tableName string) error {
	if tableExists(db, tableName) {
		return nil
	}

	// Define table creation statements
	var createStmt string
	switch tableName {
	case "leaderboards":
		createStmt = `
		CREATE TABLE "leaderboards" (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX idx_leaderboards_deleted_at ON leaderboards(deleted_at);
		`
	case "participants":
		createStmt = `
		CREATE TABLE "participants" (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX idx_participants_deleted_at ON participants(deleted_at);
		`
	case "leaderboard_entries":
		createStmt = `
		CREATE TABLE "leaderboard_entries" (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX idx_leaderboard_entries_deleted_at ON leaderboard_entries(deleted_at);
		`
	case "leaderboard_metrics":
		createStmt = `
		CREATE TABLE "leaderboard_metrics" (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX idx_leaderboard_metrics_deleted_at ON leaderboard_metrics(deleted_at);
		`
	default:
		fmt.Printf("No creation statement defined for table '%s', skipping creation\n", tableName)
		return nil
	}

	// Try to run extension if needed for UUID
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Create the table
	fmt.Printf("Creating table '%s'...\n", tableName)
	err := db.Exec(createStmt).Error
	if err != nil {
		return fmt.Errorf("error creating table %s: %w", tableName, err)
	}

	fmt.Printf("Table '%s' created successfully\n", tableName)
	return nil
}

// tableExists checks if a table exists in the database
func tableExists(db *gorm.DB, tableName string) bool {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = ?
		);
	`
	db.Raw(query, tableName).Scan(&exists)
	return exists
}

// makeColumnNullable makes a specific column in a table nullable
func makeColumnNullable(db *gorm.DB, tableName, columnName string) error {
	// First check if the column exists and has the not null constraint
	var columnInfo struct {
		ColumnName    string `gorm:"column:column_name"`
		IsNullable    string `gorm:"column:is_nullable"`
		ColumnDefault string `gorm:"column:column_default"`
		DataType      string `gorm:"column:data_type"`
	}

	result := db.Raw(`
		SELECT column_name, is_nullable, column_default, data_type  
		FROM information_schema.columns 
		WHERE table_name = ? AND column_name = ?;
	`, tableName, columnName).Scan(&columnInfo)

	if result.Error != nil {
		return fmt.Errorf("error checking column %s in table %s: %w", columnName, tableName, result.Error)
	}

	// If column exists and is not nullable, modify it
	if columnInfo.ColumnName != "" && columnInfo.IsNullable == "NO" {
		fmt.Printf("Making '%s.%s' column nullable...\n", tableName, columnName)

		// Handle specific data types differently
		switch {
		case strings.Contains(columnInfo.DataType, "char") || strings.Contains(columnInfo.DataType, "text"):
			// For string types, set empty string for nulls
			if err := db.Exec(fmt.Sprintf(`UPDATE %s SET %s = '' WHERE %s IS NULL`, tableName, columnName, columnName)).Error; err != nil {
				return fmt.Errorf("error updating null values in %s.%s: %w", tableName, columnName, err)
			}
		case strings.Contains(columnInfo.DataType, "json") || strings.Contains(columnInfo.DataType, "jsonb"):
			// For JSON types, set empty object/array
			if err := db.Exec(fmt.Sprintf(`UPDATE %s SET %s = '{}' WHERE %s IS NULL OR %s::text = 'null'`, tableName, columnName, columnName, columnName)).Error; err != nil {
				return fmt.Errorf("error updating null JSON values in %s.%s: %w", tableName, columnName, err)
			}
		case strings.Contains(columnInfo.DataType, "int"):
			// For integer fields, maybe set to 0 if appropriate
			// This depends on your business logic - be careful here
			// Consider whether a default value makes sense for your model
			if columnName == "max_entries" {
				if err := db.Exec(fmt.Sprintf(`UPDATE %s SET %s = 0 WHERE %s IS NULL`, tableName, columnName, columnName)).Error; err != nil {
					return fmt.Errorf("error updating null integer values in %s.%s: %w", tableName, columnName, err)
				}
			}
		}

		// Alter the column to allow NULL values
		if err := db.Exec(fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s DROP NOT NULL`, tableName, columnName)).Error; err != nil {
			return fmt.Errorf("error altering %s.%s column: %w", tableName, columnName, err)
		}

		fmt.Printf("Successfully made '%s.%s' column nullable\n", tableName, columnName)
	} else if columnInfo.ColumnName != "" {
		fmt.Printf("Column '%s.%s' already allows nulls\n", tableName, columnName)
	} else {
		fmt.Printf("Column '%s.%s' doesn't exist yet\n", tableName, columnName)
	}

	return nil
}

// ensureColumnWithDefault ensures a column exists and has no NULL values
// Creates the column with default value if it doesn't exist
func ensureColumnWithDefault(db *gorm.DB, tableName, columnName, defaultValue string) error {
	// Check if column exists
	var columnInfo struct {
		ColumnName    string `gorm:"column:column_name"`
		IsNullable    string `gorm:"column:is_nullable"`
		ColumnDefault string `gorm:"column:column_default"`
		DataType      string `gorm:"column:data_type"`
	}

	result := db.Raw(`
		SELECT column_name, is_nullable, column_default, data_type  
		FROM information_schema.columns 
		WHERE table_name = ? AND column_name = ?;
	`, tableName, columnName).Scan(&columnInfo)

	if result.Error != nil {
		return fmt.Errorf("error checking column %s in table %s: %w", columnName, tableName, result.Error)
	}

	// If column exists, fix any NULL values
	if columnInfo.ColumnName != "" {
		fmt.Printf("Column '%s.%s' exists, checking for NULL values...\n", tableName, columnName)

		// Check for NULL values
		var nullCount int64
		if err := db.Raw(fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s IS NULL`, tableName, columnName)).Count(&nullCount).Error; err != nil {
			return fmt.Errorf("error counting null values in %s.%s: %w", tableName, columnName, err)
		}

		if nullCount > 0 {
			fmt.Printf("Fixing %d NULL values in required column '%s.%s'...\n", nullCount, tableName, columnName)

			// Set default values for NULL fields based on data type
			var updateStmt string
			switch {
			case strings.Contains(columnInfo.DataType, "timestamp"):
				updateStmt = fmt.Sprintf(`UPDATE %s SET %s = %s WHERE %s IS NULL`,
					tableName, columnName, defaultValue, columnName)
			case strings.Contains(columnInfo.DataType, "bool"):
				updateStmt = fmt.Sprintf(`UPDATE %s SET %s = %s WHERE %s IS NULL`,
					tableName, columnName, defaultValue, columnName)
			case strings.Contains(columnInfo.DataType, "int") ||
				strings.Contains(columnInfo.DataType, "float") ||
				strings.Contains(columnInfo.DataType, "numeric") ||
				strings.Contains(columnInfo.DataType, "decimal"):
				updateStmt = fmt.Sprintf(`UPDATE %s SET %s = %s WHERE %s IS NULL`,
					tableName, columnName, defaultValue, columnName)
			default:
				// For string types
				updateStmt = fmt.Sprintf(`UPDATE %s SET %s = '%s' WHERE %s IS NULL`,
					tableName, columnName, defaultValue, columnName)
			}

			if err := db.Exec(updateStmt).Error; err != nil {
				return fmt.Errorf("error setting default values for %s.%s: %w", tableName, columnName, err)
			}

			fmt.Printf("Successfully updated NULL values in '%s.%s'\n", tableName, columnName)
		} else {
			fmt.Printf("No NULL values found in required column '%s.%s'\n", tableName, columnName)
		}
	} else {
		// Column doesn't exist, create it with the default value
		fmt.Printf("Column '%s.%s' doesn't exist, creating with default value...\n", tableName, columnName)

		// Get column type from map or use VARCHAR as fallback
		columnType, exists := fieldTypes[columnName]
		if !exists {
			columnType = "VARCHAR(255)"
			fmt.Printf("Warning: No type mapping found for column '%s', using %s\n", columnName, columnType)
		}

		// Create column with default value
		var createStmt string
		switch {
		case strings.Contains(columnType, "TIMESTAMP"):
			// For timestamp fields
			if defaultValue == "NOW()" {
				createStmt = fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s %s DEFAULT CURRENT_TIMESTAMP NOT NULL`,
					tableName, columnName, columnType)
			} else {
				createStmt = fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s %s DEFAULT %s NOT NULL`,
					tableName, columnName, columnType, defaultValue)
			}
		case strings.Contains(columnType, "BOOL") || strings.Contains(columnType, "INT") ||
			strings.Contains(columnType, "FLOAT"):
			// For numeric/boolean fields
			createStmt = fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s %s DEFAULT %s NOT NULL`,
				tableName, columnName, columnType, defaultValue)
		default:
			// For string types
			createStmt = fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s %s DEFAULT '%s' NOT NULL`,
				tableName, columnName, columnType, defaultValue)
		}

		if err := db.Exec(createStmt).Error; err != nil {
			return fmt.Errorf("error creating column %s.%s: %w", tableName, columnName, err)
		}

		fmt.Printf("Successfully created column '%s.%s' with default value\n", tableName, columnName)
	}

	return nil
}

// ensureRequiredField ensures a required field has no NULL values by setting defaults
// This is kept for backward compatibility
func ensureRequiredField(db *gorm.DB, tableName, columnName, defaultValue string) error {
	return ensureColumnWithDefault(db, tableName, columnName, defaultValue)
}

// fixBaseModelColumns ensures created_at and updated_at are properly configured
func fixBaseModelColumns(db *gorm.DB) error {
	// Get all tables
	var tables []string
	result := db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE';
	`).Scan(&tables)

	if result.Error != nil {
		return fmt.Errorf("error getting table list: %w", result.Error)
	}

	for _, table := range tables {
		// Check if created_at and updated_at columns exist
		var createdAtExists, updatedAtExists bool
		var columns []string

		db.Raw(`
			SELECT column_name 
			FROM information_schema.columns 
			WHERE table_name = ?;
		`, table).Scan(&columns)

		for _, col := range columns {
			if col == "created_at" {
				createdAtExists = true
			}
			if col == "updated_at" {
				updatedAtExists = true
			}
		}

		// Fix created_at
		if createdAtExists {
			fmt.Printf("Fixing created_at in table '%s'...\n", table)
			if err := db.Exec(fmt.Sprintf(`
				ALTER TABLE %s 
				ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
				ALTER COLUMN created_at SET NOT NULL;
			`, table)).Error; err != nil {
				fmt.Printf("Warning: Could not fix created_at in %s: %v\n", table, err)
			}
		}

		// Fix updated_at
		if updatedAtExists {
			fmt.Printf("Fixing updated_at in table '%s'...\n", table)
			if err := db.Exec(fmt.Sprintf(`
				ALTER TABLE %s 
				ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP,
				ALTER COLUMN updated_at SET NOT NULL;
			`, table)).Error; err != nil {
				fmt.Printf("Warning: Could not fix updated_at in %s: %v\n", table, err)
			}
		}
	}

	return nil
}

// Registers all migration functions
func RegisterMigrations(db *gorm.DB) error {
	if err := Migration01FixSchema(db); err != nil {
		return err
	}
	return nil
}
