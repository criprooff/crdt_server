// Package main provides a set server
package main

import (
	"log"
	"strings"

	srv "github.com/criprooff/crdt_server"
	"github.com/ianschenck/envflag"
)

type Set struct {
	items map[string]struct{}
}

func (s *Set) Add(i string) (bool, error) {
	if _, ok := s.items[i]; ok {
		return true, nil
	}
	s.items[i] = struct{}{}
	return false, nil
}

func (s *Set) Remove(i string) (bool, error) {
	if _, ok := s.items[i]; ok {
		delete(s.items, i)
		return true, nil
	}
	return false, nil
}

func (s *Set) Contains(i string) (bool, error) {
	_, ok := s.items[i]
	return ok, nil
}

func main() {
	var (
		serfAddrs    = envflag.String("SERF_CLUSTER_ADDRS", "", "host:port for serf server")
		serfBindAddr = envflag.String("SERF_AGENT_ADDR", "", "host:port for serf server")
	)
	envflag.Parse()

	addrs := strings.Split(*serfAddrs, ",")
	if len(addrs) == 0 {
		log.Fatalln("no agent addrs provided")
	}

	if err := srv.RegisterSelf(addrs, *serfBindAddr); err != nil {
		log.Fatalln("cannot register with serf:", err)
	}
	defer srv.DeregisterSelf()

	srv.AddHandler(&Set{})

	log.Println("running server...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln("cannot start server:", err)
	}
}
