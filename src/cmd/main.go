package main

import (
    "github.com/benmj87/gogo-pop3gadget/src/client"
    "github.com/benmj87/gogo-pop3gadget/src/config"
    "net"
    "crypto/tls"
)

// entry point for testing/development
func main() {
    config := config.NewConfig()
    client := client.NewClient(*config)
    client.TLSDialer = func(network string, addr string, config *tls.Config) (net.Conn, error) {
        return tls.Dial(network, addr, config)
    }

    err := client.Connect()

    if err != nil {
        panic(err)
    }

    defer client.Close()    
}