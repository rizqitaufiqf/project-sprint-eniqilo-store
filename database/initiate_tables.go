package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitiateTables(dbPool *pgxpool.Pool) error {
	// Define table creation queries
	queries := []string{
		`
        CREATE TABLE IF NOT EXISTS staffs (
            id VARCHAR(100) PRIMARY KEY NOT NULL,
            name VARCHAR(100) NOT NULL,
			phone_number VARCHAR(50) UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
        `,
		`
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(100) NOT NULL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			sku VARCHAR(50) NOT NULL,
			category VARCHAR(100) NOT NULL,
			image_url TEXT NOT NULL,
			notes VARCHAR(255) NOT NULL,
			price INT NOT NULL,
			stock INT NOT NULL,
			location VARCHAR(255) NOT NULL,
			is_available BOOL NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);		
		`,
		`
        CREATE TABLE IF NOT EXISTS customers (
            id VARCHAR(100) PRIMARY KEY NOT NULL,
            name VARCHAR(100) NOT NULL,
			phone_number VARCHAR(50) UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
        `,
		`
		CREATE TABLE IF NOT EXISTS transactions (
			id VARCHAR(100) NOT NULL PRIMARY KEY,
			customer_id VARCHAR(100) NOT NULL,
			product_details VARCHAR[] NOT NULL,
			paid INT NOT NULL,
			change INT NOT NULL,
			FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE NO ACTION,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`,
		// Add more table creation queries here if needed
	}

	// Execute table creation queries
	for _, query := range queries {
		_, err := dbPool.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	return nil
}
