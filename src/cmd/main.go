package main

import (
    "github.com/benmj87/gogo-pop3gadget/src/client"
    "github.com/benmj87/gogo-pop3gadget/src/config"
    "flag"
)

// entry point for testing/development
func main() {
    pass := flag.String("Password", "", "Password to auth with")
    username := flag.String("Username", "", "Username to auth with")
    flag.Parse()

    config := config.NewConfig()
    config.Password = *pass
    config.Username = *username

    client := client.NewClient(*config)
    err := client.Connect()
    if err != nil {
        panic(err)
    }

    defer client.Close()    

    err = client.Auth()
    if err != nil {
        panic(err)
    }

    _, _, err = client.Stat()
    if err != nil {
        panic(err)
    }

    emails, err := client.List() 
    if err != nil {
        panic(err)
    }

    for _, email := range emails {
        _, err = client.ListMessage(email.ID)
        if err != nil {
            panic(err)
        }

        _, err = client.Retrieve(email.ID) 
        if err != nil {
            panic(err)
        }

        err = client.Delete(email.ID)
        if err != nil {
            panic(err)
        }
    }

    err = client.Reset()
    if err != nil {
        panic(err)
    }
}