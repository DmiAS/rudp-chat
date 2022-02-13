package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/DmiAS/rendezvous/pkg/proto/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"chat/internal/server"
)

const (
	defaultRendezvous = "localhost:9000"
	defaultServerPort = "8080"
)

func main() {
	var clientPort string
	var rendezvous string
	var serverPort string
	var debug bool

	// default port for udp is 0 cause in this case it will be chosen automatically
	flag.BoolVar(&debug, "debug", false, "debug level for logging")
	flag.StringVar(&clientPort, "port", "0", "port to run udp")
	flag.StringVar(&rendezvous, "stun", defaultRendezvous, "rendezvous server address")
	flag.StringVar(&serverPort, "sp", defaultServerPort, "http server port")
	flag.Parse()

	// setup logger
	setupLogger(debug)

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init punch client
	cli, err := client.NewClient(ctx, clientPort, rendezvous)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// init server
	srv := server.NewServer(cli, serverPort)

	// gracefully shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Kill, os.Interrupt)
	go srv.Run()
	<-quit

	log.Info().Msg("server shutting down")
	if err := srv.Shutdown(); err != nil {
		log.Fatal().
			Err(err).
			Msg("failure to shutdown server")
	}

	log.Info().
		Msg("server stopped")
}

func setupLogger(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// to enable stack tracing
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
