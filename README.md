http-tarpit
===========

Web-server which produces infinite chunked-encoded responses. It's intended to be an active defense against malicious clients.

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
