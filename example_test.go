package gcmlib_test

import (
	"fmt"

	gcm "github.com/gamegos/gcmlib"
)

func ExampleClient_Send() {
	client := gcm.NewClient(gcm.Config{
		APIKey:     "your-gcm-api-key",
		MaxRetries: 4,
	})

	message := &gcm.Message{
		RegistrationIDs: []string{"registrationID1", "registrationID2"},
		Notification: &gcm.Notification{
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
