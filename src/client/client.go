package client

import (
	"fmt"
    "github.com/benmj87/gogo-pop3gadget/src/config"
	"crypto/tls"
    "net"
	"strings"
    "errors"
	"strconv"
)

const (
    // singleLineMessageTerminator is the standard terminator for single line commands
    // e.g. STAT
    singleLineMessageTerminator = "\r\n"
    // multiLineMessageTerminator is the standard terminator for multi-line commands
    // e.g. LIST
    multiLineMessageTerminator = ".\r\n"
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
    msg, err = c.readMsg(singleLineMessageTerminator)
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

    msg, err := c.readMsg(singleLineMessageTerminator)
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

    msg, err = c.readMsg(singleLineMessageTerminator)
    if err != nil {
        return err
    }
    if c.isError(msg) {
        return errors.New(msg)
    }

    fmt.Printf(msg)

    return nil
}

// Stat calls the stat pop3 command and returns the number of messages followed
// by the size of all the messages in bytes
func (c *Client) Stat() (uint32, uint64, error) {
    err := c.writeMsg("STAT\r\n")
    if err != nil {
        return 0, 0, err
    }

    msg, err := c.readMsg(singleLineMessageTerminator)
    if err != nil {
        return 0, 0, err
    }

    fmt.Print(msg)
    if c.isError(msg) {
        return 0, 0, errors.New(msg)
    }

    items := strings.Split(strings.Trim(msg, " \r\n\t"), " ")
    if len(items) < 3 {
        return 0, 0, fmt.Errorf("Incorrect response from STAT, not enough elements splitting on space, '%v'", msg)
    }

    totalMsgs, err := strconv.ParseUint(items[1], 10, 32)
    if err != nil {
        return 0, 0, fmt.Errorf("Incorrect message count returned %v, error was %v", items[1], err)
    }

    totalSize, err := strconv.ParseUint(items[2], 10, 64)
    if err != nil {
        return 0, 0, fmt.Errorf("Incorrect message count returned %v, error was %v", items[1], err)
    }

    return uint32(totalMsgs), totalSize, nil
}

// ListMessage calls LIST {ID} and returns the appropriate message information
func (c *Client) ListMessage(messageID int) (*Email, error) {
    err := c.writeMsg(fmt.Sprintf("LIST %v\r\n", messageID))
    if err != nil {
        return nil, err
    }

    msg, err := c.readMsg(singleLineMessageTerminator)
    if err != nil {
        return nil, err
    }

    fmt.Print(msg)

    email := NewEmail()
    err = email.ParseSingleLine(msg)
    if err != nil {
        return nil, err
    }

    return email, nil
}

// List implements the LIST call returning a list of all messages and their size
func (c *Client) List() ([]*Email, error) {
    err := c.writeMsg("LIST\r\n")
    if err != nil {
        return nil, err
    }

    msg, err := c.readMsg(multiLineMessageTerminator)
    if err != nil {
        return nil, err
    }
    
    fmt.Print(msg)

    var emails []*Email
    lines := strings.Split(msg, "\r\n")

    // remove the first item (expecting +OK) and last item (expecing terminator)
    lines = lines[1:len(lines) -2]
    for _, line := range lines {
        email := NewEmail()
        err := email.ParseLine(line)
        if err != nil {
            return nil, err
        }

        emails = append(emails, email)
    }

    return emails, nil
}

// Close issues the Quit command and closes the connection
func (c *Client) Close() error {
    defer c.connection.Close()
    err := c.writeMsg("QUIT\r\n")
    if err != nil {
        return err
    }

    msg, err := c.readMsg(singleLineMessageTerminator)
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

// readMsg reads data in chunks from the connection until terminator is detected
func (c *Client) readMsg(terminator string) (string, error) { 
    const BuffSize int = 1024

    msg := ""
    data := make([]byte, BuffSize)
    
    var err error
    var read int
    for err == nil && !strings.HasSuffix(msg, terminator) {
        read, err = c.connection.Read(data)
        msg += string(data[:read])
    }
    
    if err != nil {
        return "", err
    }

    return msg, nil
}