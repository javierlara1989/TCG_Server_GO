package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetDatabaseConfig returns database configuration from environment variables
func GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "tcg_server"),
	}
}

// Connect establishes connection to MariaDB
func Connect() error {
	config := GetDatabaseConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Successfully connected to MariaDB")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// CreateTables creates the necessary tables if they don't exist
func CreateTables() error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		validation_code VARCHAR(255) NULL,
		validation_code_expires_at TIMESTAMP NULL,
		validated_at TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createUserInfoTable := `
	CREATE TABLE IF NOT EXISTS user_info (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL UNIQUE,
		level INT NOT NULL DEFAULT 1,
		experience INT NOT NULL DEFAULT 0,
		money INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		INDEX idx_user_id (user_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createCardsTable := `
	CREATE TABLE IF NOT EXISTS cards (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type ENUM('Monster', 'Spell', 'Energy') NOT NULL,
		legend TEXT NOT NULL,
		element ENUM('Fire', 'Water', 'Wind', 'Earth', 'Neutral', 'Holy', 'Dark') NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_type (type),
		INDEX idx_element (element),
		INDEX idx_name (name)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createUserCardsTable := `
	CREATE TABLE IF NOT EXISTS user_cards (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		card_id INT NOT NULL,
		amount INT NOT NULL DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
		INDEX idx_user_id (user_id),
		INDEX idx_card_id (card_id),
		UNIQUE KEY unique_user_card (user_id, card_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createDecksTable := `
	CREATE TABLE IF NOT EXISTS decks (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		name VARCHAR(100) NOT NULL,
		valid BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		INDEX idx_user_id (user_id),
		INDEX idx_valid (valid)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createDeckCardsTable := `
	CREATE TABLE IF NOT EXISTS deck_cards (
		deck_id INT NOT NULL,
		card_id INT NOT NULL,
		number INT NOT NULL DEFAULT 1,
		PRIMARY KEY (deck_id, card_id),
		FOREIGN KEY (deck_id) REFERENCES decks(id) ON DELETE CASCADE,
		FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
		INDEX idx_deck_id (deck_id),
		INDEX idx_card_id (card_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createTablesTable := `
	CREATE TABLE IF NOT EXISTS tables (
		id INT AUTO_INCREMENT PRIMARY KEY,
		category ENUM('S','A','B','C','D') NOT NULL,
		privacy ENUM('private','public') NOT NULL,
		password VARCHAR(10) NULL,
		prize ENUM('money','card','aura') NOT NULL,
		amount INT NULL,
		winner BOOLEAN NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		finished_at TIMESTAMP NULL,
		INDEX idx_category (category),
		INDEX idx_privacy (privacy),
		INDEX idx_prize (prize),
		INDEX idx_amount (amount),
		INDEX idx_winner (winner)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createUserTablesTable := `
	CREATE TABLE IF NOT EXISTS user_tables (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		rival_id INT NULL,
		table_id INT NOT NULL,
		time INT DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (rival_id) REFERENCES users(id) ON DELETE SET NULL,
		FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE CASCADE,
		INDEX idx_user_id (user_id),
		INDEX idx_rival_id (rival_id),
		INDEX idx_table_id (table_id),
		UNIQUE KEY unique_table_user (table_id, user_id),
		UNIQUE KEY unique_table_rival (table_id, rival_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createEffectsTable := `
	CREATE TABLE IF NOT EXISTS effects (
		id INT AUTO_INCREMENT PRIMARY KEY,
		description TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL,
		INDEX idx_deleted_at (deleted_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	createCardEffectsTable := `
	CREATE TABLE IF NOT EXISTS card_effects (
		card_id INT NOT NULL,
		effect_id INT NOT NULL,
		PRIMARY KEY (card_id, effect_id),
		FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
		FOREIGN KEY (effect_id) REFERENCES effects(id) ON DELETE CASCADE,
		INDEX idx_card_id (card_id),
		INDEX idx_effect_id (effect_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	// Create users table first
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Create user_info table
	_, err = DB.Exec(createUserInfoTable)
	if err != nil {
		return fmt.Errorf("error creating user_info table: %v", err)
	}

	// Create cards table
	_, err = DB.Exec(createCardsTable)
	if err != nil {
		return fmt.Errorf("error creating cards table: %v", err)
	}

	// Create user_cards table
	_, err = DB.Exec(createUserCardsTable)
	if err != nil {
		return fmt.Errorf("error creating user_cards table: %v", err)
	}

	// Create decks table
	_, err = DB.Exec(createDecksTable)
	if err != nil {
		return fmt.Errorf("error creating decks table: %v", err)
	}

	// Create deck_cards table
	_, err = DB.Exec(createDeckCardsTable)
	if err != nil {
		return fmt.Errorf("error creating deck_cards table: %v", err)
	}

	// Create tables table
	_, err = DB.Exec(createTablesTable)
	if err != nil {
		return fmt.Errorf("error creating tables table: %v", err)
	}

	// Create user_tables table
	_, err = DB.Exec(createUserTablesTable)
	if err != nil {
		return fmt.Errorf("error creating user_tables table: %v", err)
	}

	// Create effects table
	_, err = DB.Exec(createEffectsTable)
	if err != nil {
		return fmt.Errorf("error creating effects table: %v", err)
	}

	// Create card_effects table
	_, err = DB.Exec(createCardEffectsTable)
	if err != nil {
		return fmt.Errorf("error creating card_effects table: %v", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
