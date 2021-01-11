package main

import (
	"context"
	"flag"
	"github.com/nwillc/goraft/conf"
	"github.com/nwillc/goraft/model"
	"github.com/nwillc/goraft/raftapi"
	"google.golang.org/grpc"
	"log"
)

var ClientCli struct {
	Member     *string
	ConfigFile *string
}

func SetupClientCli() {
	ClientCli.Member = flag.String("member", "none", "The member name.")
	ClientCli.ConfigFile = flag.String("config-file", conf.ConfigFile, "The configuration file.")
}

func main() {
	SetupClientCli()
	flag.Parse()
	log.Println("Start")
	config, err := model.ReadConfig(*ClientCli.ConfigFile)
	if err != nil {
		log.Fatalln("can not read config", *ClientCli.ConfigFile)
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

	response, err := api.Shutdown(ctx, &raftapi.Empty{})
	if err != nil {
		log.Fatal("Shutdown failed: ", err)
	}
	log.Println("response: ", response.Status)
	log.Println("End")
}
