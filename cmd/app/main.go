package main

import (
	"context"
	"flag"

	"github.com/DmiAS/rendezvous/pkg/proto/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	defaultRendezvous = "localhost:10000"
)

func main() {
	var clientPort string
	var rendezvous string
	var debug bool
	// default port for udp is 0 cause in this case it will be chosen automatically
	flag.BoolVar(&debug, "debug", false, "debug level for logging")
	flag.StringVar(&clientPort, "port", "0", "port to run udp")
	flag.StringVar(&rendezvous, "stun", defaultRendezvous, "rendezvous server address")

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

	// init server for static

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
