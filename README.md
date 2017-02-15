[![Build Status](https://travis-ci.org/benmj87/gogo-pop3gadget.png)](https://travis-ci.org/benmj87/gogo-pop3gadget)

# gogo-pop3gadget
Pop3 Client written in Go supporting TLS.

## Usage
For examples, see cmd/main.go where a username/password is read in from the command line and tested using gmail.

Currently, all messages sent/received will be dumped to the console for debugging and development purposes.

To create a new connection:
```
config := config.NewConfig()
config.Password = "password"
config.Username = "username"
config.Server = "mail.server.com"

client := client.NewClient(*config)
err := client.Connect()
if err != nil {
  panic(err)
}

err = client.Auth()
if err != nil {
  panic(err)
}

defer client.Close()    
```

To download all emails:
```
emails, err := client.List() 
if err != nil {
  panic(err)
}

for _, emailID := range emails {
  email, err = client.Retrieve(emailID.ID) 
  if err != nil {
    panic(err)
  }
  
  fmt.Println(email.Message)
}
```

## Configuration
Only configuration needed is:

```
type Config struct {
    // Whether to connect over TLS or not
    UseTLS bool
    // Server to connect to 
    Server string
    // Port to connect on
    Port int
    // Username to auth with
    Username string
    // Password to auth with
    Password string
}```
