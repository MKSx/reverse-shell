package main

import (
	"fmt"

	"github.com/maxlaverse/reverse-shell/util"
)

var options struct {
	CreateOptions      CreateCommand      `command:"create" description:"create a new session on a given agent"`
	AttachOptions      AttachCommand      `command:"attach" description:"attach to an existing session"`
	ListSessionOptions ListSessionCommand `command:"list-session" description:"list all the sessions available on a rendez-vous"`
	ListAgentOptions   ListAgentCommand   `command:"list-agent" description:"list all the agents available on a rendez-vous"`
	ListenOptions      ListenCommand      `command:"listen" description:"listen for agents to connect using websocket"`
	Verbose            int                `short:"v" long:"verbose" env:"VERBOSE" description:"Be verbose."`
}

type CreateCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

type AttachCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

type ListSessionCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

type ListAgentCommand struct {
	Url string `short:"U" long:"url" env:"URL" description:"Url of the rendez-vous point." required:"true"`
}

type ListenCommand struct {
	Port int `short:"P" long:"port" env:"PORT" description:"Port to listen to." required:"true"`
}

func (x *CreateCommand) Execute(args []string) error {
	sessionId := CreateSession(x.Url, args[0])
	AttachSession(x.Url, sessionId)
	return nil
}

func (x *AttachCommand) Execute(args []string) error {
	util.SetupLogging(options)
	AttachSession(x.Url, args[0])
	return nil
}

func (x *ListSessionCommand) Execute(args []string) error {
	util.SetupLogging(options)
	fmt.Printf("List of sessions:\n")
	l, err := ListSessions(x.Url)
	for _, v := range l {
		fmt.Printf(" * %s => agent: %s, masters: %s, state: %s\n", v.Name, v.Agent, v.Masters, v.State)
	}
	return err
}

func (x *ListAgentCommand) Execute(args []string) error {
	util.SetupLogging(options)
	fmt.Printf("List of agents:\n")
	l, err := ListAgents(x.Url)
	for _, v := range l {
		fmt.Printf(" * %s\n", v.Name)
	}
	return err
}

func (x *ListenCommand) Execute(args []string) error {
	util.SetupLogging(options)
	return Listen(x.Port)
}

func main() {
	util.ParseArgs(&options)
}
