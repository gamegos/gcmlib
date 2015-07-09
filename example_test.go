package gcmlib_test

import (
	"fmt"

	"github.com/gamegos/gcmlib"
)

func ExampleClient_Send() {
	client := gcmlib.NewClient("your-gcm-api-key")
	message := &gcmlib.Message{
		RegistrationIDs: []string{"registrationID1", "registrationID2"},
		Notification: &gcmlib.Notification{
			Title: "Example GCM message",
			Body:  "Hello world",
		},
		Data: map[string]string{
			"customKey": "custom value",
		},
	}

	response, err := client.Send(message)
	if err != nil {
		fmt.Printf("Error: %#v\n", err)
		return
	}

	fmt.Printf("Success: %#v\n", response)
}
