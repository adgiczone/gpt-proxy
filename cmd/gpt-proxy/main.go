package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	proxy "little-gpt/gpt-proxy/pkg"
)

var (
	printVersion = flag.Bool("version", false, "print version of this build")
	keyPath      = flag.String("keypath", "", "encry")
	certPath     = flag.String("certpath", "", "encry")
	port         = flag.Int("port", 8080, "http port")
	timeout      = flag.Int("timeout", 30, "http timeout value")
)

func main() {
	flag.Parse()
	port := *port
	timeout := *timeout
	srv, err := proxy.New(port, timeout, *certPath, *keyPath)
	if err != nil {
		log.Panic("Start proxy error", err)
	}
	go func() {
		if err := srv.Serv(); err != nil {
			log.Panic("server start error", err)
		}
	}()
	log.Info("proxy server started at:", port)
	sigs := make(chan os.Signal, 3)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	for sig := range sigs {
		log.Infof("receive signal %s, server exited", sig.String())
		os.Exit(0)
	}
}
