package main

import (
	"flag"

	"github.com/amodkala/raft"
	"github.com/amodkala/store/pkg/store"
)

func main() {

	var serverAddress, cmAddress string

	flag.StringVar(&serverAddress, "server", "localhost:8080", "the address where the server will run")
	flag.StringVar(&cmAddress, "cm", "localhost:8081", "the address where the consensus module will run")
	flag.Parse()

	cm := raft.New()
	go cm.Start(cmAddress)

	s := store.New(store.WithDB(), store.WithCM(cm))
	s.Start(serverAddress)
}
