package main

import (
	"flag"
	"fmt"
	"github.com/nwillc/goraft/conf"
	"github.com/nwillc/goraft/model"
	"log"
	"os"
)

func main() {
	conf.SetupMemberCli()
	flag.Parse()
	if *conf.MemberCli.Version {
		fmt.Printf("version %s\n", "unknown")
		os.Exit(conf.NormalExit)
	}
	log.Println("Start")

	config, err := model.ReadConfig("config.json")
	if err != nil {
		log.Fatalln("can not read config")
	}

	member, ok := config.Members[*conf.MemberCli.Member]
	if !ok {
		log.Fatalln("No config for member:", *conf.MemberCli.Member)
	}
	log.Printf("Starting member %s on port %d.\n", *conf.MemberCli.Member, member.Port)
	log.Fatalln(member.Listen(*conf.MemberCli.Member))
}

