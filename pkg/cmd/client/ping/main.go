package main

import (
	"context"
	"flag"
	"github.com/nwillc/goraft/pkg/raftapi"
	"github.com/nwillc/goraft/pkg/model"
	"google.golang.org/grpc"
	"log"
)

var ClientCli struct {
	Member *string
}

func SetupClientCli() {
	ClientCli.Member = flag.String("member", "none", "The member name.")
}

func main() {
	SetupClientCli()
	flag.Parse()
	config, err := model.ReadConfig("config.json")
	if err != nil {
		log.Fatalln("can not read config")
	}
	var member model.Member
	ok := false
	for _, m := range config.Members {
		if m.Name == *ClientCli.Member {
			ok = true
			member = m
			break
		}
	}
	if !ok {
		log.Fatalln("No config for member:", *ClientCli.Member)
	}
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(member.Address(), grpc.WithInsecure())
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
	log.Printf("Member %s: { name: %s, port: %d, role: %s }\n",
		member.String(), response.Name, response.Port, response.Role)
}
