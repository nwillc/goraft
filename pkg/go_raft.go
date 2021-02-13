package main

import (
	"flag"
	"fmt"
	"github.com/nwillc/goraft/conf"
	"github.com/nwillc/goraft/model"
	"github.com/nwillc/goraft/server"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	conf.SetupMemberCli()
	flag.Parse()
	if *conf.MemberCli.Version {
		fmt.Printf("version %s\n", "unknown")
		os.Exit(conf.NormalExit)
	}
	config, err := model.ReadConfig(*conf.MemberCli.ConfigFile)
	if err != nil {
		log.Fatalln("can not read config", *conf.MemberCli.ConfigFile)
	}

	level, err := log.ParseLevel(*conf.MemberCli.LogLevel)
	if err != nil {
		log.Fatalln("can not parse log level", *conf.MemberCli.LogLevel)
	}
	log.SetLevel(level)
	var member model.Member
	var ok = false
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
	srv := server.NewRaftServer(member, config, "")
	log.WithFields(srv.LogFields()).Infoln("Running server")
	err = srv.Run()
	if err != nil {
		log.Fatalln(err)
	}
	srv.Wait()
}
