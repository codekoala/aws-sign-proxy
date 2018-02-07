package version

import (
	"fmt"
	"runtime"
)

var (
	App       = "aws-sign-proxy"
	Version   = "v0.1.0"
	Commit    = "dev"
	BuildDate = "unknown"
)

func String() string {
	return fmt.Sprintf("%s-%s", Version, Commit)
}

func Detailed() string {
	return fmt.Sprintf(`%s %s
Commit:		%s
Build date:	%s
Go:		%s`,
		App, Version, Commit, BuildDate, runtime.Version())
}
