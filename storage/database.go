package storage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Storage abstraction
type Database struct {
	DB                    *sqlx.DB // sqlx database connection pool
	host                  string
	port                  int
	user                  string
	password              string
	name                  string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
	log                   *zap.Logger
}

// NewDatabaseOptions is the options for creating a new database connection
type NewDatabaseOptions struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
	Log                   *zap.Logger
}

// Constructor for Database
func NewDatabase(opts NewDatabaseOptions) *Database {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	return &Database{
		host:                  opts.Host,
		port:                  opts.Port,
		user:                  opts.User,
		password:              opts.Password,
		name:                  opts.Name,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
		log:                   opts.Log,
	}
}

// Connect to database
func (d *Database) Connect() error {
	d.log.Info("Connecting to database", zap.String("url", d.createDatasourceName(false)))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	d.DB, err = sqlx.ConnectContext(ctx, "pgx", d.createDatasourceName(true))
	if err != nil {
		return err
	}

	d.log.Debug("Setting database connection pool settings",
		zap.Int("max open connections", d.maxOpenConnections),
		zap.Int("max idle connections", d.maxIdleConnections),
		zap.Duration("connection max life", d.connectionMaxLifetime),
		zap.Duration("connection max idle time", d.connectionMaxIdleTime))
	d.DB.SetMaxOpenConns(d.maxOpenConnections)
	d.DB.SetMaxIdleConns(d.maxIdleConnections)
	d.DB.SetConnMaxLifetime(d.connectionMaxLifetime)
	d.DB.SetConnMaxIdleTime(d.connectionMaxIdleTime)

	return nil

}

// Builds the connection URL for the database; allows for password to be hidden or not
func (d *Database) createDatasourceName(withPassword bool) string {
	password := d.password
	if !withPassword {
		password = "xxx"
	}

	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", d.user, password, d.host, d.port, d.name)
}

// Ping database to check if it is alive
func (d *Database) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := d.DB.PingContext(ctx); err != nil {
		return err
	}

	_, err := d.DB.ExecContext(ctx, "SELECT 1") // Test query
	return err
}
