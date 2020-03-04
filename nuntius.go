package main

import (
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"flag"
	"fmt"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type PushMessage struct {
	BasicNotification *messaging.Notification `json:"notification,omitempty"`
	Data              map[string]string       `json:"data,omitempty"`
}

func main() {
	var (
		credentialsFile = flag.String("credentialsFile", "",
			`A Firebase credentials file downloaded from the Firebase console.
Log into the Firebase console and go to your project. 
Beside "Project Overview" on the left click the gear/settings icon.  
Select "Project Settings".  In the Project Settings page, click on "Service accounts".  
Scroll down and click on "Generate new private key".  This is your credentials file.`)
		topic    = flag.String("topic", "", "topic to send the message to")
		token    = flag.String("token", "", "token to send the message to")
		pushFile = flag.String("pushFile", "",
			`A file that contains a notification (title, body, image (optional)) 
and/or data (a map of key/value pairs) that will be pushed to the client(s). 
Format:
{
   "notification" : { 
      "title" : "Test Title",
      "body"  : "Test Body",
      "image" : "https://whatever.com/image.png"
   },  

   "data" : { 
     "key1" : "value1",
     "key2" : "value2"
   }   
}`)

		ttl = flag.Int64("ttl", 0, "Time-to-live value for notifications in seconds.\n"+
			"0 (default) means \"now or never\", that is,\n"+
			"deliver the message now or don't deliver it at all.\n"+
			"Max value is 2419200 (28 days).\n"+
			"For details see: https://firebase.google.com/docs/cloud-messaging/concept-options")
	)

	flag.Parse()

	if *credentialsFile == "" {
		log.Fatalln("credentialsFile not specified")
	}

	if *pushFile == "" {
		log.Fatalln("pushFile not specified")
	}

	if *topic == "" && *token == "" {
		log.Fatalf("A token or topic must be specified")
	}

	if *topic != "" && *token != "" {
		log.Fatalf("You can specify only a token or topic, not both.")
	}

	// https://firebase.google.com/docs/cloud-messaging/concept-options
	if *ttl < 0 || *ttl > 2419200 {
		log.Fatalf("ttl must be between 0 and 2419200 (28 days) inclusive.  0 means deliver immediately and don't retry on failure.")
	}

	var firebaseApp = initializeAppWithServiceAccount(*credentialsFile)

	// https://github.com/firebase/snippets-go/blob/master/admin/messaging.go
	ctx := context.Background()

	client, err := firebaseApp.Messaging(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	pushMsg, err := readFileAsPushMessage(*pushFile)

	if err != nil {
		log.Fatalln(err)
	}

	if pushMsg.BasicNotification != nil {
		if pushMsg.BasicNotification.Title == "" && pushMsg.BasicNotification.Body != "" {
			log.Fatal("In your notification file, if body is set, then title must also be set")
		}

		if pushMsg.BasicNotification.Title != "" && pushMsg.BasicNotification.Body == "" {
			log.Fatal("In your notification file, if title is set, then body must also be set")
		}
	} else {
		if len(pushMsg.Data) == 0 {
			log.Fatal("The pushFile must contains either a notification or a non-empty map of data or both.")
		}
	}

	var androidTimeToLive time.Duration

	if *ttl != 0 {
		androidTimeToLive = time.Duration(*ttl) * time.Second
	}

	APNS_headers := map[string]string{
		"apns-expiration": fmt.Sprintf("0"),
	}

	UTC, err := time.LoadLocation("UTC")

	if err == nil && *ttl > 0 {
		APNS_headers["apns-expiration"] = fmt.Sprintf("%v", time.Now().In(UTC).Unix()+*ttl)
	}

	Webpush_headers := map[string]string{
		"TTL": fmt.Sprintf("%v", *ttl),
	}

	message := &messaging.Message{
		Notification: pushMsg.BasicNotification,
		Data:         pushMsg.Data,
		Android: &messaging.AndroidConfig{
			TTL: &androidTimeToLive,
		},
		APNS: &messaging.APNSConfig{
			Headers: APNS_headers,
		},
		Webpush: &messaging.WebpushConfig{
			Headers: Webpush_headers,
		},
	}

	if *topic != "" {
		message.Topic = *topic
	} else {
		message.Token = *token
	}

	response, err := client.Send(ctx, message)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Successfully sent message:", response)
}

// https://github.com/firebase/snippets-go/blob/master/admin/main.go
func initializeAppWithServiceAccount(credentialsFile string) *firebase.App {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	return app
}

func readFileAsPushMessage(fileName string) (*PushMessage, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening file: %v", err))
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading file: %v", err))
	}

	pushMsg := new(PushMessage)

	err = json.Unmarshal(bytes, pushMsg)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing file: %v", err))
	}

	return pushMsg, nil
}
