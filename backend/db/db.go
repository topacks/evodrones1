package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/root/evodrones/backend/config"
)

type DB struct {
	*sqlx.DB
}

var db *DB

func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	db = &DB{conn}
	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetDB retorna a instância do banco de dados
func GetDB() *DB {
	return db
}

// Executar transação
func WithTransaction(fn func(*sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
