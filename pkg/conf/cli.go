package conf

import "flag"

// MemberCli are the CLI flags for member operations.
var MemberCli struct {
	Member   *string
	Version  *bool
	LogLevel *string
}

// SetupMemberCli the CLI flags for members.
func SetupMemberCli() {
	MemberCli.Member = flag.String("member", "one", "The member name.")
	MemberCli.Version = flag.Bool("version", false, "Display version.")
	MemberCli.LogLevel = flag.String("log-level", "info", "Logging level.")
}
