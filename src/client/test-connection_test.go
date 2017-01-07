package client

import (
    "testing"
    "errors"
    "time"
)

// Test_TestConnectionReadOk tests that when read is called the correct data is returned
func Test_TestConnectionReadOk(t *testing.T) {
    toTest := NewTestConnection()
    toTest.ToRead = append(toTest.ToRead, "FOO1\r\n")

    buff := make([]byte, 3)
    read, err := toTest.Read(buff) 
    if read != 3 || err != nil {
        t.Errorf("Incorrect read %v or error %v", read, err)
    }
    if string(buff) != "FOO" {
        t.Errorf("Incorrect data read, was %v", string(buff))
    }
    
    read, err = toTest.Read(buff)
    if read != 3 || err != nil {
        t.Errorf("Incorrect read %v or error %v", read, err)
    }
    if string(buff) != "1\r\n" {
        t.Errorf("Incorrect data read, was %v", string(buff))
    }
}

// Test_TestConnectionReadErrors checks various errors are returned
func Test_TestConnectionReadErrors(t *testing.T) {
    toTest := NewTestConnection()  
    
    read, err := toTest.Read([]byte{})
    if read != 0 || err != nil {
        t.Error("Unexpected data read")
    }

    toTest.ReadError = errors.New("test")
    read, err = toTest.Read([]byte{})
    if read != 0 || err == nil {
        t.Error("Expected an error to be returned")
    }
    
    read, err = toTest.Read(nil)
    if read != 0 || err == nil {
        t.Error("Expected an error to be returned")
    }
}

// Test_LocalAddr tests an error is returned
func Test_LocalAddr(t *testing.T) {
    toTest := NewTestConnection()  
    ret := toTest.LocalAddr()
    if ret != nil {
        t.Error("Expected nil to be returned")
    }
}

// Test_RemoteAddr tests an error is returned
func Test_RemoteAddr(t *testing.T) {
    toTest := NewTestConnection()  
    ret := toTest.RemoteAddr()
    if ret != nil {
        t.Error("Expected nil to be returned")
    }
}

// Test_SetDeadline tests an error is returned
func Test_SetDeadline(t *testing.T) {
    toTest := NewTestConnection()  
    ret := toTest.SetDeadline(time.Now())
    if ret == nil {
        t.Error("Expected an error to be returned")
    }
}

// Test_SetReadDeadline tests an error is returned
func Test_SetReadDeadline(t *testing.T) {
    toTest := NewTestConnection()  
    ret := toTest.SetReadDeadline(time.Now())
    if ret == nil {
        t.Error("Expected an error to be returned")
    }
}

// Test_SetWriteDeadline tests an error is returned
func Test_SetWriteDeadline(t *testing.T) {
    toTest := NewTestConnection()  
    ret := toTest.SetWriteDeadline(time.Now())
    if ret == nil {
        t.Error("Expected an error to be returned")
    }
}

// Test_Close Tests the close functionality
func Test_Close(t *testing.T) {
    toTest := NewTestConnection()

    err := toTest.Close()
    if !toTest.Closed || err != nil {
        t.Error("Expected closed to be set")
    }

    toTest.CloseError = errors.New("fo")
    toTest.Closed = false
    err = toTest.Close()
    if toTest.Closed || err == nil {
        t.Error("Expected an error to be returned")
    }
}

// Test_Write tests the write functionality
func Test_Write(t *testing.T) {
    toTest := NewTestConnection()
    written, err := toTest.Write([]byte("a"))
    if written != 1 || err != nil || len(toTest.Written) == 0 || toTest.Written[0] != "a" {
        t.Errorf("Unexpected response when writing, written %v, error %v, msg %v", written, err, toTest.Written)
    }

    toTest = NewTestConnection()
    toTest.WriteError = errors.New("fo")
    written, err = toTest.Write([]byte { 'a' })
    if written != 0 || err == nil {
        t.Errorf("Unexpected response when writing")
    }
    
    toTest = NewTestConnection()
    toTest.WriteCount = 0
    written, err = toTest.Write([]byte { 'a' })
    if written != 0 || err != nil {
        t.Errorf("Unexpected response when writing")
    }
}