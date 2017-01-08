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
    testConn.ThrowReadErrorAfter = 1
    err := toTest.Close()
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_AuthOk tests that auth is ok
func Test_AuthOk(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = false
    conf.Password = "p4ssw0rd"
    conf.Username = "foo@bar.com"

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.Dialer = func(net string, server string) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    
    testConn.ToRead = append(testConn.ToRead, "+OK Username ok\r\n")
    testConn.ToRead = append(testConn.ToRead, "+OK Pass ok\r\n")
    err := toTest.Auth()

    if err != nil {
        t.Error("Error returned")
    }
    if testConn.Written[0] != "USER foo@bar.com\r\n" {
        t.Error("Invalid USER written")
    }
    if testConn.Written[1] != "PASS p4ssw0rd\r\n" {
        t.Error("Invalid password written")
    }
}

// Test_AuthUsernameError tests that an error is returned on read
func Test_AuthUsernameError(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = false
    conf.Password = "p4ssw0rd"
    conf.Username = "foo@bar.com"

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.Dialer = func(net string, server string) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    
    testConn.WriteError = errors.New("foo")
    testConn.TimesWriteCalled = 0
    testConn.TimesReadCalled = 0
    err := toTest.Auth()

    if err == nil {
        t.Error("Error not returned")
    }

    testConn.WriteError = nil
    testConn.TimesWriteCalled = 0
    testConn.TimesReadCalled = 0
    testConn.ReadError = errors.New("foo")
    err = toTest.Auth()

    if err == nil {
        t.Error("Error not returned")
    }
    
    testConn.WriteError = nil
    testConn.ReadError = nil
    testConn.TimesWriteCalled = 0
    testConn.TimesReadCalled = 0
    testConn.ToRead = append(testConn.ToRead, "-ERR sumthing wrong\r\n")
    err = toTest.Auth()

    if err == nil {
        t.Error("Error not returned")
    }
}

// Test_AuthPasswordWriteError tests that an error is returned on writing the password
func Test_AuthPasswordWriteError(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = false
    conf.Password = "p4ssw0rd"
    conf.Username = "foo@bar.com"

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.Dialer = func(net string, server string) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.WriteError = errors.New("foo")
    testConn.ThrowWriteErrorAfter = 1 // first write is USER
    err := toTest.Auth()

    if err == nil {
        t.Error("Error not returned")
    }
}

// Test_AuthPasswordReadError tests that an error is returned on writing the password
func Test_AuthPasswordReadError(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = false
    conf.Password = "p4ssw0rd"
    conf.Username = "foo@bar.com"

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.Dialer = func(net string, server string) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.ReadError = errors.New("foo")
    testConn.ThrowReadErrorAfter = 2 // includes connecting and USER
    err := toTest.Auth()

    if err == nil {
        t.Error("Error not returned")
    }
}

// Test_AuthPasswordReadErrorMsgReturned tests that an error is returned on reading the response from the password
func Test_AuthPasswordReadErrorMsgReturned(t *testing.T) {
    conf := config.NewConfig()
    conf.UseTLS = false
    conf.Password = "p4ssw0rd"
    conf.Username = "foo@bar.com"

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.Dialer = func(net string, server string) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")
    testConn.ToRead = append(testConn.ToRead, "-ER\r\n")
    err := toTest.Auth()

    if err == nil {
        t.Error("Error not returned")
    }
}