package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	"github.com/nwillc/goraft/model"
	"github.com/nwillc/goraft/setup"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type server struct {
	raftapi.UnimplementedRaftServiceServer
}

func main() {
	flag.Parse()
	if *setup.Flags.Version {
		fmt.Printf("version %s\n", "unknown")
		os.Exit(setup.NormalExit)
	}
	log.Println("Start")

	config, err := model.ReadConfig("config.json")
	if err != nil {
		log.Fatalln("can not read config")
	}

	member, ok := config.Members[*setup.Flags.Member]
	if !ok {
		log.Fatalln("No config for member:", *setup.Flags.Member)
	}
	log.Printf("Starting member %s on port %d.\n", *setup.Flags.Member, member.Port)
	listen, err := net.Listen("tcp", member.Address())
	if err != nil {
		log.Fatalln(err)
	}
	srv := grpc.NewServer()
	raftapi.RegisterRaftServiceServer(srv, &server{})
	log.Fatalln(srv.Serve(listen))
}

func (s *server) Ping(ctx context.Context, request *raftapi.Empty) (*raftapi.Bool, error) {
	return &raftapi.Bool{Status: true}, nil
}
