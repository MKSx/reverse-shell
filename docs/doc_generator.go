package main

import (
	"fmt"

	agents "github.com/maxlaverse/reverse-shell/agent/cmd"
	client "github.com/maxlaverse/reverse-shell/client/cmd"
	rendezvous "github.com/maxlaverse/reverse-shell/rendezvous/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	generate(agents.GetCommand(), "./docs/agent")
	generate(client.GetCommand(), "./docs/client")
	generate(rendezvous.GetCommand(), "./docs/rendezvous")
}

func generate(command *cobra.Command, dest string) {
	err := doc.GenMarkdownTree(command, dest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Documentation successfully generated")
}
