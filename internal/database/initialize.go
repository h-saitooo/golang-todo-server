package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// DB is the global database connection pool
	DB *sql.DB
)

func Initialize() error {
	var err error
	jst, err := time.LoadLocation("Asia/Tokyo")

	if err != nil {
		return fmt.Errorf("failed to load location: %w", err)
	}

	fmt.Printf("MYSQL_DATABASE: %s\n", os.Getenv("MYSQL_DATABASE"))
	fmt.Printf("MYSQL_USER: %s\n", os.Getenv("MYSQL_USER"))
	fmt.Printf("MYSQL_PASSWORD: %s\n", os.Getenv("MYSQL_PASSWORD"))

	config := mysql.Config{
		DBName:    os.Getenv("MYSQL_DATABASE"),
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Net:       "tcp",
		Addr:      "database:3306",
		Collation: "utf8mb4_general_ci",
		Loc:       jst,
	}

	DB, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool parameters
	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(10)
	DB.SetConnMaxLifetime(time.Hour)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	// Create table query
	query := `
		CREATE TABLE IF NOT EXISTS users
		(
			id int(11) NOT NULL AUTO_INCREMENT,
			name varchar(50) NOT NULL,
			email varchar(255) NOT NULL,
			password_hash varchar(255) NOT NULL,
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			PRIMARY KEY (id)
		)
		ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`

	// Exec query to mysql database
	_, err = DB.Exec(query)

	fmt.Println("Successfully created table: users")

	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
