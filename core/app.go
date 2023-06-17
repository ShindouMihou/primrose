package main

import (
	"github.com/joho/godotenv"
	"github.com/kataras/golog"
	"os"
	"os/signal"
	"primrose/clients"
	"primrose/routes"
	"primrose/toolkit"
	"sync"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		golog.New().SetPrefix("[ENV] ").Fatal("Failed to load environment files: ", err)
	}
	clients.M.Init()
	if isMainApp := toolkit.HandleArguments(); isMainApp {
		routes.Conduct()
		clients.Iris.Init()
		<-Shutdown(5 * time.Second)
	}
}

func Shutdown(timeout time.Duration) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		out := time.AfterFunc(timeout, func() {
			golog.Fatal("Failed to graceful shutdown after timeout, forcing exit.")
		})
		defer out.Stop()
		var wg sync.WaitGroup
		wg.Wait()
		close(wait)
	}()
	return wait
}
