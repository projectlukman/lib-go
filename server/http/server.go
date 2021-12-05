package http

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/projectlukman/lib-go/log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

type ConfigServer struct {
	ServiceName      string
	Address          string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	GracefullTimeout time.Duration
}

type serverHandler struct {
	configServer ConfigServer
	router       *router.Router
}

func NewHttpServer(configServer ConfigServer, router *router.Router) Server {
	return &serverHandler{
		configServer: configServer,
		router:       router,
	}
}

func (s *serverHandler) Run(ctx context.Context) error {
	server := fasthttp.Server{
		ReadTimeout:  s.configServer.ReadTimeout,
		WriteTimeout: s.configServer.WriteTimeout,
		Name:         s.configServer.ServiceName,
		Handler:      s.router.Handler,
	}

	// NOTE: Package reuseport provides a TCP net.Listener with SO_REUSEPORT support.
	// SO_REUSEPORT allows linear scaling server performance on multi-CPU servers.

	// create a fast listener ;)
	ln, err := reuseport.Listen("tcp4", s.configServer.Address)
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
		return err
	}

	// create a graceful shutdown listener
	graceful := NewGracefulListener(ln, s.configServer.GracefullTimeout)

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname unavailable: %s", err)
		return err
	}

	// Error handling
	listenErr := make(chan error, 1)

	/// Run server
	go func() {
		log.Printf("%s - Web server starting on %v", hostname, graceful.Addr())
		log.Printf("%s - Press Ctrl+C to stop", hostname)
		// listenErr <- s.HTTPServer.ListenAndServe(":" + cfg.Port)
		listenErr <- server.Serve(graceful)
	}()

	// SIGINT/SIGTERM handling
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Handle channels/graceful shutdown
	for {
		select {
		// If server.ListenAndServe() cannot start due to errors such
		// as "port in use" it will return an error.
		case err := <-listenErr:
			if err != nil {
				log.Fatalf("listener error: %s", err)
			}
			os.Exit(0)
		// handle termination signal
		case <-osSignals:
			fmt.Printf("\n")
			log.Printf("%s - Shutdown signal received.\n", hostname)

			// Servers in the process of shutting down should disable KeepAlives
			// FIXME: This causes a data race
			server.DisableKeepalive = true

			// Attempt the graceful shutdown by closing the listener
			// and completing all inflight requests.
			if err := graceful.Close(); err != nil {
				log.Fatalf("error with graceful close: %s", err)
				return err
			}

			log.Printf("%s - Server gracefully stopped.\n", hostname)
		}
	}
}
