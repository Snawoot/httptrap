package main

import (
    "time"
    "net/http"
    "bytes"
)

const TARGET_BUF = 128 * 1024

type SlowResponder struct {
    interval time.Duration
    contentType string
    content []byte
    logger *CondLogger
}

func NewSlowResponder(interval time.Duration,
                      content []byte,
                      contentType string,
                      logger *CondLogger) *SlowResponder {
    if contentType == "" {
        contentType = "text/html"
    }
    return &SlowResponder{
        interval: interval,
        contentType: contentType,
        content: content,
        logger: logger,
    }
}

func (s *SlowResponder) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
    s.logger.Info("Client %s connected", req.RemoteAddr)
    defer s.logger.Info("Client %s disconnected", req.RemoteAddr)
    wr.Header().Set("Content-Type", s.contentType)
    wr.WriteHeader(http.StatusOK)
    flusher, flusherOk := wr.(http.Flusher)
    ctx := req.Context()
    if ! flusherOk {
        s.logger.Critical("Server doesn't support response flushing!")
        return
    }
    for {
        _, err := wr.Write(s.content)
        if err != nil {
            break
        }
        flusher.Flush()
        select {
        case <-time.After(s.interval):
        case <-ctx.Done():
            return
        }
    }
}

type FastResponder struct {
    contentType string
    content []byte
    logger *CondLogger
}

func NewFastResponder(content []byte,
                      contentType string,
                      logger *CondLogger) *SlowResponder {
    if contentType == "" {
        contentType = "text/html"
    }
    buf := bytes.Repeat(content, TARGET_BUF / len(content))
    return &SlowResponder{
        contentType: contentType,
        content: buf,
        logger: logger,
    }
}

func (s *FastResponder) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
    s.logger.Info("Client %s connected", req.RemoteAddr)
    defer s.logger.Info("Client %s disconnected", req.RemoteAddr)
    wr.Header().Set("Content-Type", s.contentType)
    wr.WriteHeader(http.StatusOK)
    for {
        _, err := wr.Write(s.content)
        if err != nil {
            break
        }
    }
}

