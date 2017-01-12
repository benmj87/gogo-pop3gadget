package client

import (
    "testing"
)

// Email_ParseLineOk checks that a line has been parsed correctly
func Test_EmailParseLineOk(t *testing.T) {
    toTest := NewEmail()

    err := toTest.ParseLine("10 10")    
    if err != nil {
        t.Error(err)
    }

    if toTest.Size != 10 {
        t.Error("Incorrect size")
    }
    if toTest.ID != 10 {
        t.Error("Incorrect ID")
    }
}

// Test_EmailParseLineErrorsReturned checks that an error is returned correctly
func Test_EmailParseLineErrorsReturned(t *testing.T) {
    toTest := NewEmail()

    err := toTest.ParseLine("10")    
    if err == nil {
        t.Error("Expected error")
    }

    err = toTest.ParseLine("a 10")    
    if err == nil {
        t.Error("Expected error")
    }

    err = toTest.ParseLine("10 a")    
    if err == nil {
        t.Error("Expected error")
    }
}

// Test_EmailParseSingleLineOk checks that a line has been parsed correctly
func Test_EmailParseSingleLineOk(t *testing.T) {
    toTest := NewEmail()

    err := toTest.ParseSingleLine("+OK 10 10")    
    if err != nil {
        t.Error(err)
    }

    if toTest.Size != 10 {
        t.Error("Incorrect size")
    }
    if toTest.ID != 10 {
        t.Error("Incorrect ID")
    }
}

// Test_EmailParseSingleLineErrorsReturned checks that an error is returned correctly
func Test_EmailParseSingleLineErrorsReturned(t *testing.T) {
    toTest := NewEmail()

    err := toTest.ParseSingleLine("+OK 10")    
    if err == nil {
        t.Error("Expected error")
    }

    err = toTest.ParseSingleLine("+OK a 10")    
    if err == nil {
        t.Error("Expected error")
    }

    err = toTest.ParseSingleLine("+OK 10 a")    
    if err == nil {
        t.Error("Expected error")
    }
}