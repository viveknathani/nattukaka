package database

import (
	"log"
	"os"
	"testing"

	"github.com/viveknathani/nattukaka/config"
)

func init() {
	config.LoadEnvFromFile("../.env")
}
func TestInitializeAndClose(t *testing.T) {
	db := &Database{}
	dsn := os.Getenv("DATABASE_URL")
	err := db.Initialize(dsn)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.pool.Query("select version();")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
