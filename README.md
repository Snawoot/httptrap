httptrap
========

Web-server which produces infinite chunked-encoded responses. It's intended to be a part of active defense against malicious HTTP clients to cancel their load impact and/or bruteforce efforts. Depending on settings, httptrap may keep busy most of attacker's resources, diverting them to wait slowly fed infinite response, or cause often OOM and botnet payload crash if httptrap sends response as fast as possible and client uses RAM to buffer entire response.

## Why?

Why just not collect attackers IP addresses and ban them? Or maybe even automate this process and establish it on ongoing basis?

Filtering approach has some drawbacks:

* It prevents access of legitimate users who share IP address with attacker (e.g. common ISP NAT gateway, compromised devices in network, common proxy or TOR exit node and so on). Consequently, it gives attacker an additional power to cut someone's access from service if one will manage originate flood traffic from same IP. Slowdown approach is more gentle and capable to not impact correct requests from legitimate users.
* Filtering may have lower stopping power if attacker posesses a huge pool of source addresses. In fact, if pool is sufficiently large and it's full rotation period is longer than IP cooldown time, attackers may conduct their activities continously, without downtimes.

## Installation

Pre-built binaries available on [releases](https://github.com/Snawoot/httptrap/releases/latest) page.

Alternatively, you may install httptrap from source:

```
go get github.com/Snawoot/httptrap
```

## Synopsis

```
$ ./httptrap -h
Usage of ./httptrap:
  -bind string
    	listen address (default ":8008")
  -cert string
    	enable HTTPS and use certificate
  -ct string
    	Content-Type value for responses (default "text/html")
  -interval duration
    	interval between chunks (default 1s)
  -key string
    	key for TLS certificate
  -string string
    	hex-encoded representation of byte string repeated in responses (default "0a")
  -verbosity int
    	logging verbosity (10 - debug, 20 - info, 30 - warning, 40 - error, 50 - critical) (default 20)
```
