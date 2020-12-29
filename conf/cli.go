package conf

import "flag"

var MemberCli struct {
	Member  *string
	Version *bool
}

func SetupMemberCli() {
	MemberCli.Member = flag.String("member", "one", "The member name.")
	MemberCli.Version = flag.Bool("version", false, "Display version.")
}
