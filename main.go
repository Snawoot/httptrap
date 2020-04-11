package main

import (
    "fmt"
    "time"
    "flag"
    "os"
    "net/http"
    "log"
    "encoding/hex"
)

func perror(msg string) {
    fmt.Fprintln(os.Stderr, "")
    fmt.Fprintln(os.Stderr, msg)
}

func arg_fail(msg string) {
    perror(msg)
    perror("Usage:")
    flag.PrintDefaults()
    os.Exit(2)
}

type CLIArgs struct {
    bind string
    verbosity int
    interval time.Duration
    cert, key string
    contentType string
    content []byte
}


func parse_args() *CLIArgs {
    var (
        args CLIArgs
        content string
        err error
    )
    flag.StringVar(&args.bind, "bind", ":8008", "listen address")
    flag.IntVar(&args.verbosity, "verbosity", 20, "logging verbosity " +
                "(10 - debug, 20 - info, 30 - warning, 40 - error, 50 - critical)")
    flag.DurationVar(&args.interval, "interval", time.Second, "interval between chunks")
    flag.StringVar(&args.cert, "cert", "", "enable HTTPS and use certificate")
    flag.StringVar(&args.key, "key", "", "key for TLS certificate")
    flag.StringVar(&args.contentType, "ct", "text/html", "Content-Type value for responses")
    flag.StringVar(&content, "string", "0a", "hex-encoded string repeated in responses")
    flag.Parse()
    args.content, err = hex.DecodeString(content)
    if err != nil {
        arg_fail(err.Error())
    }
    return &args
}

func main() {
    var (
        err error
    )
    args := parse_args()

    logWriter := NewLogWriter(os.Stderr)
    defer logWriter.Close()

    mainLogger := NewCondLogger(log.New(logWriter, "MAIN    : ",
                                log.LstdFlags | log.Lshortfile),
                                args.verbosity)
    handlerLogger := NewCondLogger(log.New(logWriter, "HANDLER : ",
                                   log.LstdFlags | log.Lshortfile),
                                   args.verbosity)
    mainLogger.Info("Starting server...")

    var server http.Server
    server.Addr = args.bind
    var handler http.Handler
    if args.interval == 0 {
        handler = NewFastResponder(args.content, args.contentType, handlerLogger)
    } else {
        handler = NewSlowResponder(args.interval, args.content, args.contentType, handlerLogger)
    }
    server.Handler = handler
    server.ErrorLog = log.New(logWriter, "HTTPSRV : ", log.LstdFlags | log.Lshortfile)
    if args.cert != "" {
        cfg, err := makeServerTLSConfig(args.cert, args.key, "")
        if err != nil {
            mainLogger.Critical("TLS config construction failed: %v", err)
            os.Exit(3)
        }
        server.TLSConfig = cfg
        err = server.ListenAndServeTLS("", "")
    } else {
        err = server.ListenAndServe()
    }

    mainLogger.Critical("Server terminated with a reason: %v", err)
    mainLogger.Info("Shutting down...")
}
