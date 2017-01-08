package client

import (
    "time"
    "net"
    "errors"
)

// TestConnection is used for testing connections
type TestConnection struct {
    // ToRead holds a list of items to be read in order
    ToRead []string
    // Written holds a list of strings written appended
    Written []string
    // ReadError holds an error to return on Read
    ReadError error
    // WriteError holds an error to return on Read
    WriteError error
    // CloseError holds an error to return upon Close
    CloseError error
    // WriteCount allows the length to be set when calling write
    WriteCount int
    // Closed holds whether close has been called
    Closed bool
    // TimesReadCalled holds the number of times read has been called
    TimesReadCalled int
    // TimesWriteCalled holds the number of times Write has been called
    TimesWriteCalled int
    // ThrowReadErrorAfter is the TimesWriteCalled to throw the error after
    ThrowReadErrorAfter int
    // ThrowWriteErrorAfter is the TimesWriteCalled to throw the error after
    ThrowWriteErrorAfter int
}

// NewTestConnection returns a new TestConnection
func NewTestConnection() *TestConnection {
    return &TestConnection {
        ToRead: make([]string, 0),
        Written: make([]string, 0),
        WriteCount: -1,
        Closed: false,
        TimesReadCalled: 0,
        TimesWriteCalled: 0,
        ThrowWriteErrorAfter: 0,
        ThrowReadErrorAfter: 0,
    }
}

// Reads data from the ToRead array
func (c *TestConnection) Read(b []byte) (n int, err error) {
    toRet := 0
    if b == nil {
        return 0, errors.New("b cannot be nil")
    }

    if c.ReadError != nil && c.TimesReadCalled == c.ThrowReadErrorAfter {
        return 0, c.ReadError
    }

    if len(c.ToRead) == 0 {
        return 0, nil
    } 
        
    dataToRet := c.ToRead[0]
    buffLength := len(b)
    
    // b is big enough to hold dataToRet
    if buffLength >= len(dataToRet) {
        copy(b, []byte(dataToRet))
        c.ToRead = append(c.ToRead[:0], c.ToRead[1:]...) // remove the first element 
        toRet = len(dataToRet)
    } else {
        // need to only return the maximum we can
        remains := dataToRet[buffLength:len(dataToRet)]
        c.ToRead[0] = remains // keep the remainder of the data
        copy(b, dataToRet[0:buffLength])
        toRet = buffLength
    }
    
    c.TimesReadCalled++
    return toRet, nil
}

// Write writes data to Written.
func (c *TestConnection) Write(b []byte) (n int, err error) {
    if c.WriteError != nil && c.ThrowWriteErrorAfter == c.TimesWriteCalled {
        return 0, c.WriteError
    }

    if c.WriteCount > -1 {
        return c.WriteCount, nil
    }

    c.TimesWriteCalled++
    c.Written = append(c.Written, string(b))
    return len(b), nil
}

// Close closes the connection.
func (c *TestConnection) Close() error {
    if c.CloseError != nil {
        return c.CloseError
    }
    
    c.Closed = true
    return nil
}

// LocalAddr returns the local network address.
func (c *TestConnection) LocalAddr() net.Addr {
    return nil
}

// RemoteAddr returns the remote network address.
func (c *TestConnection) RemoteAddr() net. Addr {
    return nil
}

// SetDeadline not yet implemented
func (c *TestConnection) SetDeadline(t time.Time) error {
    return errors.New("Not implemented")
}

// SetReadDeadline not yet implemented
func (c *TestConnection) SetReadDeadline(t time.Time) error {
    return errors.New("Not implemented")
}

// SetWriteDeadline not yet implemented
func (c *TestConnection) SetWriteDeadline(t time.Time) error {
    return errors.New("Not implemented")
}