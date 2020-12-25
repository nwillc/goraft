package main

import (
	"context"
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	"github.com/nwillc/goraft/model"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	raftapi.UnimplementedRaftServiceServer
}

func main() {
	log.Println("Start")

	config, err := model.ReadConfig("config.json")
	if err != nil {
		log.Fatalln("can not read config")
	}

	port := fmt.Sprintf(":%d", config.Members[0].Port)
	listen, err := net.Listen("tcp", port)
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
