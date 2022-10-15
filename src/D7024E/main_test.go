package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestArgs(t *testing.T) {
	/*
		fmt.Println("implement main testing")

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		cases := []struct {
			Empty   string
			Command string
			Address string
			Port    string
			Testing string
		}{
			{".", "start", "", "8080", "test"},

			{".", "join", "172.17.0.2", "8080", "test"},
			{".", "ping", "172.17.0.2", "8080", "test"},

			{".", "ping", "172.17.1.30", "8080", "test"},

			//{".", "get", "", "HashValue", "test"},
			//{".", "put", "", "HashValue", "test"},
		}
		for _, tc := range cases {

			// this call is required because otherwise flags panics, if args are set between flag.Parse calls
			flag.CommandLine = flag.NewFlagSet(tc.Command, flag.ExitOnError)
			// we need a value to set Args[0] to, cause flag begins parsing at Args[1]
			os.Args = []string{tc.Empty}

			if tc.Command != "" {
				os.Args = append(os.Args, tc.Command)
			}

			if tc.Address != "" {
				os.Args = append(os.Args, tc.Address)
			}
			if tc.Port != "" {
				os.Args = append(os.Args, tc.Port)
			}
			if tc.Testing != "" {
				os.Args = append(os.Args, tc.Testing)
			}

			main()

		}

		fmt.Println("-------------------------")
		fmt.Println("")
	*/
}
