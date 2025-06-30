package database

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	migrate "github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"

	// For golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// For golang-migrate
	_ "github.com/lib/pq"
)

// NewDatabase use gorm to connect to our database and returns the connection
func NewDatabase(url string) *gorm.DB {
	connection, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: connection,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("Error using existing connection:", err)
	}

	return db
}

// RunMigrations will run migrations from the migrations directory on your database connecction.
func RunMigrations(relativePathToDirectory string, connection *sql.DB) {

	// Initialize the database driver for migrations
	driver, err := migratePostgres.WithInstance(connection, &migratePostgres.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver: %v", err)
	}

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance("file://"+relativePathToDirectory, "postgres", driver)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("migrations applied successfully!")
}
