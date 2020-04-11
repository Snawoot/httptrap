package main

import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "errors"
)

func makeServerTLSConfig(certfile, keyfile, cafile string) (*tls.Config, error) {
    var cfg tls.Config
    cert, err := tls.LoadX509KeyPair(certfile, keyfile)
    if err != nil {
        return nil, err
    }
    cfg.Certificates = []tls.Certificate{cert}
    if cafile != "" {
        roots := x509.NewCertPool()
        certs, err := ioutil.ReadFile(cafile)
        if err != nil {
            return nil, err
        }
        if ok := roots.AppendCertsFromPEM(certs); !ok {
            return nil, errors.New("Failed to load CA certificates")
        }
        cfg.ClientCAs = roots
        cfg.ClientAuth = tls.VerifyClientCertIfGiven
    }
    return &cfg, nil
}
