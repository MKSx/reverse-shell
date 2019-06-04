package main

import "github.com/maxlaverse/reverse-shell/rendezvous/cmd"

func main() {

	if err := cmd.GetCommand().Execute(); err != nil {
		panic(err)
	}
}
