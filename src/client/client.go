package client

import (
	"fmt"
    "github.com/benmj87/gogo-pop3gadget/src/config"
	"crypto/tls"
    "net"
	"strings"
    "errors"
)

// Client holds code for the connection
type Client struct {
    // the config for the connection
    config config.Config
    // the connection
    connection net.Conn
    // the dialer
    Dialer func(string, string) (net.Conn, error)
    // the tls dialer to create new tls connections
    TLSDialer func(string, string, *tls.Config) (net.Conn, error)
}

// NewClient returns a new default instance of the Client
func NewClient(conf config.Config) *Client {
    return &Client{
        config: conf,
        Dialer: net.Dial,
    }
}

// Connect opens the connection and initiates
func (c *Client) Connect() error {
    var err error

    if (c.config.UseTLS) {
        fmt.Printf("Connecting using TLS to %v:%v\n", c.config.Server, c.config.Port)
        c.connection, err = c.TLSDialer("tcp", fmt.Sprintf("%v:%v", c.config.Server, c.config.Port), &tls.Config{})
    } else {
        fmt.Printf("Connecting to %v:%v\n", c.config.Server, c.config.Port)
        c.connection, err = c.Dialer("tcp", c.config.Server + ":" + string(c.config.Port))
    }
    
    if err != nil {
        return err
    }

    msg := ""
    msg, err = c.readMsg()
    if err != nil {
        return err
    }
    
    if c.isError(msg) {
        return errors.New(msg)
    }
    
    fmt.Print(msg)

    return nil
}

// Auth calls USER + PASS
func (c *Client) Auth() error {
    err := c.writeMsg(fmt.Sprintf("USER %v\r\n", c.config.Username))
    if err != nil {
        return err
    }

    msg, err := c.readMsg()
    if err != nil {
        return err
    }
    if c.isError(msg) {
        return errors.New(msg)
    }

    fmt.Printf(msg)

    err = c.writeMsg(fmt.Sprintf("PASS %v\r\n", c.config.Password))
    if err != nil {
        return err
    }

    msg, err = c.readMsg()
    if err != nil {
        return err
    }
    if c.isError(msg) {
        return errors.New(msg)
    }

    fmt.Printf(msg)

    return nil
}

// Close issues the Quit command and closes the connection
func (c *Client) Close() error {
    defer c.connection.Close()
    err := c.writeMsg("QUIT\r\n")
    if err != nil {
        return err
    }

    msg, err := c.readMsg()
    if err != nil {
        return err
    }

    fmt.Print(msg)

    return nil
}

// isError checks if the string starts with -ERR or !+OK
func (c *Client) isError(msg string) bool {
    if strings.HasPrefix(msg, "-ERR") {
        return true
    } 
    if !strings.HasPrefix(msg, "+OK") {
        return true
    }

    return false
}

// writeMsg writes the data to the connection and checks for errors
func (c *Client) writeMsg(msg string) error {
    fmt.Print(msg)
    written, err := c.connection.Write([]byte(msg))

    if err != nil {
        return err
    }

    if written != len(msg) {
        return fmt.Errorf("Invalid length of data written to connection, expected %v but only managed %v", len(msg), written)
    }

    return nil
}

// readMsg reads data in chunks from the connection until \r\n is detected
func (c *Client) readMsg() (string, error) { 
    const BuffSize int = 1024
       
    msg := ""
    data := make([]byte, BuffSize)
    
    var err error
    var read int
    for err == nil && !strings.HasSuffix(msg, "\r\n") {
        read, err = c.connection.Read(data)
        msg += string(data[:read])
    }
    
    if err != nil {
        return "", err
    }

    return msg, nil
}