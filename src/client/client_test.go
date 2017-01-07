package client

import (
    "testing"
    "github.com/benmj87/gogo-pop3gadget/src/config"
	"crypto/tls"
	"net"
    "errors"
)

// Test_TLSConnectOk tests that connect works ok for TLS
func Test_TLSConnectOk(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    err := toTest.Connect()
    if err != nil {
        t.Error("Error returned")
    }
}

// Test_NonTLSConnectOk tests that connect works ok for NonTLS
func Test_NonTLSConnectOk(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = false

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.Dialer = func(net string, server string) (net.Conn, error) {
        return testConn, nil
    }
    
    err := toTest.Connect()
    if err != nil {
        t.Error("Error returned")
    }
}

// Test_TLSConnectReturnsError tests that an error is returned
func Test_TLSConnectReturnsError(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, errors.New("foo")
    }
    
    err := toTest.Connect()
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_TLSConnectErrorOnRead tests that an error is returned on read
func Test_TLSConnectErrorOnRead(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.ReadError = errors.New("foo")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    err := toTest.Connect()
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_TLSConnectErrorMsgReturned tests that an error is returned on read
func Test_TLSConnectErrorMsgReturned(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "-ERR\r\n")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    err := toTest.Connect()
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_TLSConnectNonSuccessMsgReturned tests that an error is returned on read
func Test_TLSConnectNonSuccessMsgReturned(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "ERM\r\n")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    err := toTest.Connect()
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_CloseOk checks close is called
func Test_CloseOk(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.ToRead = append(testConn.ToRead, "+GREAT\r\n")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    err := toTest.Close()
    if err != nil {
        t.Error("Error returned")
    }

    if !testConn.Closed {
        t.Error("Close wasn't called")
    }

    if testConn.Written[0] != "QUIT\r\n" {
        t.Error("QUIT wasn't written")
    }
}

// Test_CloseErrorReturned checks an error is returned correctly from close
func Test_CloseErrorReturned(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.ToRead = append(testConn.ToRead, "+GREAT\r\n")
    testConn.WriteError = errors.New("foo")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    err := toTest.Close()    
    if err == nil {
        t.Error("Error not returned")
    }
}


// Test_CloseErrorOnWrite checks an error is returned on write
func Test_CloseErrorOnWrite(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.ToRead = append(testConn.ToRead, "+GREAT\r\n")
    testConn.WriteError = errors.New("foo")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    err := toTest.Close()
    if err == nil {
        t.Error("Error not returned")
    }
}

// Test_CloseIncorrectWriteLength returns error when write returns a different length
func Test_CloseIncorrectWriteLength(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.ToRead = append(testConn.ToRead, "+GREAT\r\n")
    testConn.WriteCount = 2 // not the same length of the QUIT msg

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    err := toTest.Close()
    if err == nil {
        t.Error("Error not returned")
    }
}

// Test_CloseErrorOnRead checks error is returned on read
func Test_CloseErrorOnRead(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = true

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.ToRead = append(testConn.ToRead, "+GREAT\r\n")

    toTest := NewClient(*conf)
    toTest.TLSDialer = func(net string, server string, tlsConf *tls.Config) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    
    testConn.ReadError = errors.New("foo")
    err := toTest.Close()
    if err == nil {
        t.Error("No error returned")
    }
}