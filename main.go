package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/e-dard/lumen/service"
)

func main() {
	v := flag.Bool("v", false, "Turns on verbose logging.")
	d := flag.Duration("refresh", 10*time.Second, "Refresh duration.")
	flag.Parse()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	options := []service.Option{service.WithInterval(*d)}
	if *v {
		options = append(options, service.Verbose())
	}

	s := service.NewService("/tmp/lumen.db", options...)
	s.Start()

	<-quit
	s.Close()
}
