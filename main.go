package main

import (
	"context"
	"github.com/nwillc/goraft/api/raftapi"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

type server struct {
	raftapi.UnimplementedRaftServiceServer
}

func main() {
	log.Println("Start", port)

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	srv := grpc.NewServer()
	raftapi.RegisterRaftServiceServer(srv, &server{})
	log.Fatalln(srv.Serve(listen))
	log.Println("Stop")
}

func (s *server) Ping(ctx context.Context, request *raftapi.Empty) (*raftapi.Bool, error) {
	return &raftapi.Bool{Status: true}, nil
}
