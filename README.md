# Netassertv2-l4-client

The `Netassertv2-l4-client` is a Go program designed to open TCP / UDP connections toward a specified destination and send a configurable string payload.

A test comprises one or more connection attempts, and a passed test results in the client exiting with a status code of 0 (1 otherwise).

You can pull the latest Docker image from `docker.io/controlplane/netassertv2-l4-client:1.0.0`

## Configuration

The client accepts the following environment variables:

| Environment Variable | Go Type | Default Value | Purpose |
| --- | --- | --- | --- |
| TARGET_HOST | string |  | The target host to test (mandatory), it can be an IP or a domain name |
| TARGET_PORT | uint16 |  | The target port to test (mandatory) |
| PROTOCOL | string | tcp | Protocol, tcp or udp |
| MESSAGE | string | defaultmessage | The string to use as payload |
| PERIOD | uint | 5000 | Interval between connections in milliseconds |
| TIMEOUT | uint | 2000 | Period of time in milliseconds after which the connection is considered failed |
| ATTEMPTS | uint | 1 | Number of connections to create (regardless of their failed / successful statuses) |
| SUCCESS_THRESHOLD | uint | 80 | Percentage of successful connections above which the whole test is considered passed |
| LOG_LEVEL | string | info | Log level, available: (debug, info, warn, error) |
| LOG_ENCODING | string | console | Log encoding, available: (console, json) |

Parameters can also be passed as command line arguments. These take precedence over environment variables in case both are defined for the same parameter. Command line args are named as lowered-case environment variables and with "`-`" replaced by "`_`".

## Build and Run

Build the client:

```bash
make build
```

And run it providing at least a target host and a port, for example:

```bash
./bin/netassertv2-l4-client --target-host 192.168.1.10 --target-port 8443
```

If the host is not reachable and / or there is not a service listening on the specified port, the connection(s) will fail.

Multiple connections can run in parallel, for example when connections fail due to timeout and the `period` parameter is less than `timeout`, or, more in general, if a connection takes more than `period` to complete.

## Example

A valid host can be created using netcat in server mode:

```bash
while true; do nc -vl localhost 12345; done
```

On a different tab, run the client:

```bash
./bin/netassertv2-l4-client --target-host=localhost --target-port=12345 --attempts 3 --message $'examplemessage\n'
```

Client output:

```bash
2023/03/09 15:49:38 maxprocs: Leaving GOMAXPROCS=16: CPU quota undefined
2023-03-09T15:49:38.166+0100	INFO	cobra@v1.6.1/command.go:916	&{LogLevel:info LogEncoding:console Protocol:tcp TargetHost:localhost TargetPort:12345 Message:examplemessage
 Timeout:2000 Attempts:3 Period:5000 SuccThrPec:80}
2023-03-09T15:49:38.168+0100	INFO	runtime/asm_amd64.s:1598	successful connection and data sent to localhost:12345
2023-03-09T15:49:43.171+0100	INFO	runtime/asm_amd64.s:1598	successful connection and data sent to localhost:12345
2023-03-09T15:49:48.170+0100	INFO	runtime/asm_amd64.s:1598	done creating connections
2023-03-09T15:49:48.171+0100	INFO	conntester/conntester.go:93	waiting for connections to stop...
2023-03-09T15:49:48.172+0100	INFO	runtime/asm_amd64.s:1598	successful connection and data sent to localhost:12345
2023-03-09T15:49:48.172+0100	INFO	conntester/conntester.go:93	all connections have finished
2023-03-09T15:49:48.173+0100	INFO	runtime/asm_amd64.s:1598	success rate of: 100
2023-03-09T15:49:48.173+0100	INFO	runtime/asm_amd64.s:1598	success rate greater than threshold: 80
2023-03-09T15:49:48.173+0100	INFO	cobra@v1.6.1/command.go:916	jobs stopped
2023-03-09T15:49:48.173+0100	INFO	cobra@v1.6.1/command.go:916	test passed
```

Netcat output showing the "`examplemessage`" received:

```bash
Listening on localhost 12345
Connection received on localhost 33226
examplemessage
Listening on localhost 12345
Connection received on localhost 33228
examplemessage
Listening on localhost 12345
Connection received on localhost 55004
examplemessage
Listening on localhost 12345
```