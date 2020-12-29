package main

import (
	"context"
	"flag"
	"github.com/nwillc/goraft/api/raftapi"
	"github.com/nwillc/goraft/model"
	"google.golang.org/grpc"
	"log"
)

var ClientCli struct {
	Member *string
}

func SetupClientCli()  {
	ClientCli.Member = flag.String("member", "none", "The member name.")
}

func main() {
	SetupClientCli()
	flag.Parse()
	log.Println("Start")
	config, err := model.ReadConfig("config.json")
		if err != nil {
			log.Fatalln("can not read config")
		}
	member, ok := config.Members[*ClientCli.Member]
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
	log.Println("response: ", response.Name, ":", response.Port)
	log.Println("End")
}
