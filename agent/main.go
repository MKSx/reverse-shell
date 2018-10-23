package main

import (
	"github.com/maxlaverse/reverse-shell/agent/cmd"
)

func main() {
	if err := cmd.GetCommand().Execute(); err != nil {
		panic(err)
	}
}
