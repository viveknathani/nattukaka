package cache

import (
	"log"
	"testing"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAndSet(t *testing.T) {

	url := "127.0.0.1:6379"
	c := &Cache{}
	c.Initialize(url, "", "")
	conn := c.Pool.Get()

	_, err := Set(conn, "random", []byte("45"))
	handleError(err)

	// Get a value that exists
	value, err := Get(conn, "random")
	handleError(err)
	if string(value) != "45" {
		log.Fatal("Incorrect GET")
	}

	// Get a value that does not exist
	dne, err := Get(conn, "dne")
	handleError(err)
	if dne != nil {
		log.Fatal("Getting something that does not exist.")
	}

	err = conn.Close()
	handleError(err)

	err = c.Pool.Close()
	handleError(err)
}
