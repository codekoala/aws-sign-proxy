package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codekoala/aws-sign-proxy/version"
)

func init() {
	var showVersion bool

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `aws-sign-proxy %s

aws-sign-proxy allows you to easily sign requests using AWS v4 Signatures
without having to update your actual application code with signing logic.
Simply issue the same request that you would normally send through
aws-sign-proxy, and the request will be signed before being sent to the target
service.

Usage:
        %s

For configuration options, see the README.

Options:
`, version.Version, os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.Parse()

	if showVersion {
		fmt.Println(version.Detailed())
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}
