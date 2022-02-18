
# Firebase Cloud Messaging

Golang client library for Firebase Cloud Messaging.
It uses Firebase Cloud Messaging HTTP protocol: https://firebase.google.com/docs/cloud-messaging/http-server-ref.

## Getting Started

```sh
$ go get github.com/dulabs/go-fcm
```

## Send message

```go

package main

import (
	"log"

	"github.com/dulabs/go-fcm"
)

func main() {
    // Create a client
    client, err := fcm.NewClient("your-api-key", "your-server-key")

	if err != nil {
		log.Fatalln(err)
	}

    // New message
    message := &fcm.Message{
        To: "device-token-or-topic", // /topics/topic-name
        Notification: &fcm.Notification{
            Title: "Title",
            Body:  "Body",
        },
    }

    // Send the message
    resp, err := client.Send(message)
    if err != nil {
        log.Fatalln(err)
    }

    log.Printf("%#v\n", resp)
}
```
## Subscribe

```go
// Subscribe
err := client.SubscribeToTopic([]string{"device-token"}, "/topics/test-topic-name")

// Unsubscribe
err = client.UnsubscribeFromTopic([]string{"device-token"}, "/topics/test-topic-name")
```

## Device Info

```go
// Get device info
info, err := client.GetDeviceInfo("device-token")
if err != nil {
    fmt.Println(info)
}
```
