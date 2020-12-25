package main

import (
	"context"
	"github.com/nwillc/goraft/api/raftapi"
	"google.golang.org/grpc"
	"log"
)

func main() {
	log.Println("Start")
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect", err)
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	ctx := context.Background()

	response, err := api.Ping(ctx, &raftapi.Empty{})
	if err != nil {
		log.Fatal("Ping failed: ", err)
	}
	log.Println("response: ", response.Status)
	log.Println("End")
}
