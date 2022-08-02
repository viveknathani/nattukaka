package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/viveknathani/nattukaka/server"
)

var port string = ""

func init() {
	port = os.Getenv("PORT")
}

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.NewServer()
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
	fmt.Println("Goodbye!")
}
