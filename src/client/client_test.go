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

// Test_StatOk checks that Stat functions correctly
func Test_StatOk(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK 10 1024 vunderbar\r\n")

    msgs, size, err := toTest.Stat()

    if err != nil {
        t.Error("Error returned")
    }
    
    if testConn.Written[0] != "STAT\r\n" {
        t.Error("Stat written incorrectly") 
    }
    
    if msgs != 10 || size != 1024 {
        t.Error("Invalid msg count or size returned")
    }
}

// Test_StatErrorOnReadWrite checks that errors are returned correctly on ReadWrite
func Test_StatErrorOnReadWrite(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()
    
    testConn.ReadError = nil
    testConn.WriteError = errors.New("Foo")
    _, _, err := toTest.Stat()
    if err == nil {
        t.Error("No error returned")
    }

    testConn.ReadError = errors.New("Foo")
    testConn.WriteError = nil
    _, _, err = toTest.Stat()
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_StatInvalidMsg checks that errors are handled with incorrect responses returned
func Test_StatInvalidMsg(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "-ERR uh oh\r\n")    
    _, _, err := toTest.Stat()
    if err == nil {
        t.Error("No error returned")
    }
    
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")    
    _, _, err = toTest.Stat()
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_StatInvalidInt checks that errors are handled with incorrect ints returned
func Test_StatInvalidInt(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "+OK -1 a\r\n")    
    _, _, err := toTest.Stat()
    if err == nil {
        t.Error("No error returned")
    }
    
    testConn.ToRead = append(testConn.ToRead, "+OK 10 a\r\n")    
    _, _, err = toTest.Stat()
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_ListOk checks that a listing works correctly
func Test_ListOk(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "+OK\r\n1 10\r\n2 4\r\n.\r\n")
    emails, err := toTest.List()
    if err != nil {
        t.Error(err)
    }

    if len(emails) != 2 {
        t.Errorf("Invalid length %v", len(emails))
    }
    if emails[0].ID != 1 || emails[1].ID != 2 {
        t.Error("Invalid emails parsed")
    }
}

// Test_ListReadWriteError checks that a read and write error returns correctly
func Test_ListReadWriteError(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.WriteError = errors.New("foo")
    _, err := toTest.List()
    if err == nil {
        t.Error("Expected an error")
    }
    
    testConn.WriteError = nil
    testConn.ReadError = errors.New("foo")
    _, err = toTest.List()
    if err == nil {
        t.Error("Expected an error")
    }
}

// Test_ListInvalidData checks that an invalid response returns correctly
func Test_ListInvalidData(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "+OK\r\na\r\n.\r\n")
    _, err := toTest.List()
    if err == nil {
        t.Error("Expected an error")
    }
}

// Test_ListMessageOk checks that List is called correctly
func Test_ListMessageOk(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "+OK 10 100\r\n")
    email, err := toTest.ListMessage(10)
    if err != nil {
        t.Error("An error happened")
    }

    if email.ID != 10 || email.Size != 100 {
        t.Error("Invalid ID or Size")
    }
}

// Test_ListMessageReadWriteError checks that a read and write error returns correctly
func Test_ListMessageReadWriteError(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.WriteError = errors.New("foo")
    _, err := toTest.ListMessage(10)
    if err == nil {
        t.Error("Expected an error")
    }
    
    testConn.WriteError = nil
    testConn.ReadError = errors.New("foo")
    _, err = toTest.ListMessage(10)
    if err == nil {
        t.Error("Expected an error")
    }
}

// Test_ListMessageInvalidData checks that an invalid response returns correctly
func Test_ListMessageInvalidData(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "+OK 10\r\n")
    _, err := toTest.ListMessage(10)
    if err == nil {
        t.Error("Expected an error")
    }
    
    testConn.ToRead = append(testConn.ToRead, "-ERR\r\n")
    _, err = toTest.ListMessage(10)
    if err == nil {
        t.Error("Expected an error")
    }
}

// Test_RetrieveOk Checks that a message is retrieved correctly
func Test_RetrieveOk(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "+OK 100\r\nThe message\r\n.\r\n")
    email, err := toTest.Retrieve(10)

    if err != nil {
        t.Error("Error returned")
    }
    if testConn.Written[0] != "RETR 10\r\n" {
        t.Error("Invalid command")
    }
    if email.ID != 10 && email.Message != "The message" {
        t.Error("Invalid message")
    }
}

// Test_RetrieveErrorReturned Checks that a message error is returned
func Test_RetrieveErrorReturned(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "-ERR 100\r\nThe message\r\n.\r\n")
    _, err := toTest.Retrieve(10)
    if err == nil {
        t.Error("No error returned")
    }
}

// Test_RetrieveReadWriteError checks for read and write errors
func Test_RetrieveReadWriteError(t *testing.T) {
testConn, toTest, _ := initialiseConnection()
    testConn.WriteError = errors.New("foo")
    _, err := toTest.Retrieve(10)
    if err == nil {
        t.Error("Expected an error")
    }
    
    testConn.WriteError = nil
    testConn.ReadError = errors.New("foo")
    _, err = toTest.Retrieve(10)
    if err == nil {
        t.Error("Expected an error")
    }
}

// Test_DeleteOk checks that DELE is called correctly
func Test_DeleteOk(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "+OK deleted\r\n")
    err := toTest.Delete(10)

    if err != nil {
        t.Error("Error returned")
    }
    if testConn.Written[0] != "DELE 10\r\n" {
        t.Error("Invalid command")
    }
}

// Test_DeleteReadWriteError checks that when DELE is called a read write error
func Test_DeleteReadWriteError(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "+OK deleted\r\n")
    testConn.WriteError = errors.New("foo")
    testConn.ReadError = nil
    err := toTest.Delete(10)

    if err == nil {
        t.Error("No error returned")
    }
    
    testConn.ToRead = append(testConn.ToRead, "+OK deleted\r\n")
    testConn.ReadError = errors.New("foo")
    testConn.WriteError = nil
    err = toTest.Delete(10)

    if err == nil {
        t.Error("No error returned")
    }
}

// Test_DeleteErrorMsg checks that when DELE is called and an error thats returned is handled
func Test_DeleteErrorMsg(t *testing.T) {
    testConn, toTest, _ := initialiseConnection()

    testConn.ToRead = append(testConn.ToRead, "-ERR deleted\r\n")
    err := toTest.Delete(10)

    if err == nil {
        t.Error("No error returned")
    }
}

// initialiseConnection initialises a connection calling connect
// and resetting any read counters back to 0
func initialiseConnection() (*TestConnection, *Client, *config.Config) {
    conf := config.NewConfig()
    conf.UseTLS = false

    testConn := NewTestConnection()
    testConn.ToRead = append(testConn.ToRead, "+OK\r\n")

    toTest := NewClient(*conf)
    toTest.Dialer = func(net string, server string) (net.Conn, error) {
        return testConn, nil
    }
    
    toTest.Connect()
    testConn.TimesReadCalled = 0
    
    return testConn, toTest, conf
}