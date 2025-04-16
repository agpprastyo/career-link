package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// PostgresDB wraps database functionality
type PostgresDB struct {
	db  *sql.DB
	log *logger.Logger
	mu  sync.Mutex
}

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(cfg config.DatabaseConfig, log *logger.Logger) (*PostgresDB, error) {
	log.Info("Establishing PostgreSQL connection")
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.WithError(err).Error("Failed to open database connection")
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	log.Info("Configuring connection pool settings")
	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Testing database connection")
	if err := db.PingContext(ctx); err != nil {
		log.WithError(err).Error("Failed to ping database")
		err := db.Close()
		if err != nil {
			log.WithError(err).Error("Failed to close database connection after ping error")
			return nil, err
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL database")
	return &PostgresDB{db: db, log: log}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	p.log.Info("Closing database connection")
	return p.db.Close()
}

// GetDB returns the underlying database connection
func (p *PostgresDB) GetDB() *sql.DB {
	return p.db
}

// Ping checks database connection
func (p *PostgresDB) Ping(ctx context.Context) error {
	p.log.Debug("Pinging database")
	return p.db.PingContext(ctx)
}

// ExecContext executes a query without returning any rows
func (p *PostgresDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	p.log.WithField("query", query).Debug("Executing SQL statement")
	result, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		p.log.WithError(err).WithField("query", query).Error("SQL execution failed")
	}
	return result, err
}

// QueryContext executes a query that returns rows
func (p *PostgresDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	p.log.WithField("query", query).Debug("Executing SQL query")
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		p.log.WithError(err).WithField("query", query).Error("SQL query failed")
	}
	return rows, err
}

// QueryRowContext executes a query that returns at most one row
func (p *PostgresDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	p.log.WithField("query", query).Debug("Executing SQL query row")
	return p.db.QueryRowContext(ctx, query, args...)
}

// BeginTx starts a transaction
func (p *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	p.log.Debug("Starting database transaction")
	tx, err := p.db.BeginTx(ctx, opts)
	if err != nil {
		p.log.WithError(err).Error("Failed to begin transaction")
	}
	return tx, err
}

func (p *PostgresDB) GetContext(ctx context.Context, i *int, query string) interface{} {
	p.log.WithField("query", query).Debug("Executing get context query")
	return p.db.QueryRowContext(ctx, query).Scan(i)
}

func (p *PostgresDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	p.log.WithField("query", query).Debug("Executing select context query")
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		p.log.WithError(err).WithField("query", query).Error("Select query failed")
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			p.log.WithError(err).Error("Failed to close rows")
		}
	}(rows)
	return sqlx.StructScan(rows, dest)
}
