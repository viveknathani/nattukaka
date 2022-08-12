package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/viveknathani/nattukaka/cache"
	"github.com/viveknathani/nattukaka/database"
	"github.com/viveknathani/nattukaka/server"
	"github.com/viveknathani/nattukaka/service"
)

var port string = ""
var databaseServer string = ""
var redisServer string = ""

func init() {
	port = os.Getenv("PORT")
	databaseServer = os.Getenv("DATABASE_URL")
	redisServer = os.Getenv("REDIS_URL")
}

// getDatabase will init and return a db
func getDatabase() *database.Database {

	db := &database.Database{}
	err := db.Initialize(databaseServer)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return db
}

// getCache will return a connection to Redis from the pool
func getCache() (*cache.Cache, redis.Conn) {

	memory := &cache.Cache{}
	memory.Initialize(redisServer, "", "")
	memoryConn := memory.Pool.Get()
	return memory, memoryConn
}

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	db := getDatabase()
	memory, memoryConn := getCache()
	srv := &server.Server{
		Service: &service.Service{
			Repo: db,
			Conn: memoryConn,
		},
		Router: mux.NewRouter(),
	}
	srv.SetupRoutes()
	go func() {
		err := http.ListenAndServe(":"+port, srv)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}()
	fmt.Println("Server started!")
	<-done
	shutdown(srv, db, memory)
}

func shutdown(srv *server.Server, db *database.Database, memory *cache.Cache) {

	err := srv.Service.Conn.Close()
	if err != nil {
		fmt.Print(err)
	}
	err = memory.Close()
	if err != nil {
		fmt.Print(err)
	}
	err = db.Close()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("goodbye!")
}
