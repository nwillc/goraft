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
	Value      *int
}

func SetupClientCli() {
	ClientCli.Member = flag.String("member", "none", "The member name.")
	ClientCli.ConfigFile = flag.String("config-file", conf.ConfigFile, "The configuration file.")
	ClientCli.Value = flag.Int("value", 42, "Value to append to log.")
}

func main() {
	SetupClientCli()
	flag.Parse()
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

	response, err := api.AppendValue(ctx, &raftapi.Value{Value: int64(*ClientCli.Value)})
	if err != nil {
		log.Fatal("Append failed: ", err)
	}
	log.Println("Response", response)
}
