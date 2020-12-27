package setup

import "flag"

var Flags struct {
	Member  *string
	Version *bool
}

func init() {
	Flags.Member = flag.String("member", "one", "The member name.")
	Flags.Version = flag.Bool("version", false, "Display version.")
}
