package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	"github.com/nwillc/goraft/model"
	"log"
	"os"
)

type server struct {
	raftapi.UnimplementedRaftServiceServer
}

func main() {
	SetupMemberCli()
	flag.Parse()
	if *MemberCli.Version {
		fmt.Printf("version %s\n", "unknown")
		os.Exit(NormalExit)
	}
	log.Println("Start")

	config, err := model.ReadConfig("config.json")
	if err != nil {
		log.Fatalln("can not read config")
	}

	member, ok := config.Members[*MemberCli.Member]
	if !ok {
		log.Fatalln("No config for member:", *MemberCli.Member)
	}
	log.Printf("Starting member %s on port %d.\n", *MemberCli.Member, member.Port)
	log.Fatalln(member.Listen())
}

func (s *server) Ping(ctx context.Context, request *raftapi.Empty) (*raftapi.Bool, error) {
	return &raftapi.Bool{Status: true}, nil
}
