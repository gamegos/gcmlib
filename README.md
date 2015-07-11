# gcmlib

[![Build Status](https://travis-ci.org/gamegos/gcmlib.svg?branch=master)](https://travis-ci.org/gamegos/gcmlib)

Golang Google Cloud Messaging(GCM) library.


## Installation
```
$ go get github.com/gamegos/gcmlib
```

Then

```go
import "github.com/gamegos/gcmlib"
```


## Example Usage

```go
client := gcmlib.NewClient(&gcmlib.Options{APIKey: "your-gcm-api-key"})
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

```


## Tests

`gcmlib` provides both unit and integration tests.

To run unit tests:

```
$ go test
```

To run integration tests you need to pass your android application's `GCM` **api key** and a **registration id** for this application.

```
$ go test -v -tags=integration -key=$GCM_KEY -regid=$GCM_REGID
```

By default, push messages in integration tests will only be sent to the google servers in `dry_run` mode. If you actually want to deliver messages to the device, set ```-dry=false```

```
$ go test -v -tags=integration -key=$GCM_KEY -regid=$GCM_REGID -dry=false
```



## License

MIT. See [LICENSE](./LICENSE).
