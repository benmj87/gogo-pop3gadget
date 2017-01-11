package client

import (
    "strings"
    "fmt"
    "strconv"
)

// Email holds information relating to a single email
type Email struct {
    // ID holds the message id
    ID int
    // Size holds the message size in bytes
    Size uint
}

// NewEmail creates a new email based upon the ID and Size
func NewEmail() *Email {
    return &Email{}
}

// ParseLine parses a line expecting +OK {ID} {SIZE}
func (e *Email) ParseLine(line string) error {
    items := strings.Split(strings.Trim(line, " \r\n\t"), " ")
    if len(items) < 2 {
        return fmt.Errorf("Incorrect line, not enough elements splitting on space, '%v'", line)
    }

    id, err := strconv.ParseInt(items[0], 10, 32)
    if err != nil {
        return fmt.Errorf("Incorrect message count returned %v, error was %v", items[1], line)
    }

    totalSize, err := strconv.ParseUint(items[1], 10, 32)
    if err != nil {
        return fmt.Errorf("Incorrect message count returned %v, error was %v", items[1], err)
    }

    e.Size = uint(totalSize)
    e.ID = int(id)

    return nil
}