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

	var member model.Member
	ok := false
	for _, m := range config.Members {
		if m.Name == *conf.MemberCli.Member {
			ok = true
			member = m
			break
		}
	}
	if !ok {
		log.Fatalln("No config for member:", *conf.MemberCli.Member)
	}
	srv :=model.NewServer(member, config)
	log.Fatalln(srv.Run())
}
