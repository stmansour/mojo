package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mojo/util"
	"net/http"
	"strings"
	"time"
)

// handles subscription messages from AWS SNS

// The following HTTP POST request is an example of a subscription confirmation message.
//
// POST / HTTP/1.1
// x-amz-sns-message-type: SubscriptionConfirmation
// x-amz-sns-message-id: 165545c9-2a5c-472c-8df2-7ff2be2b3b1b
// x-amz-sns-topic-arn: arn:aws:sns:us-west-2:123456789012:MyTopic
// Content-Length: 1336
// Content-Type: text/plain; charset=UTF-8
// Host: example.com
// Connection: Keep-Alive
// User-Agent: Amazon Simple Notification Service Agent
// {
//   "Type" : "SubscriptionConfirmation",
//   "MessageId" : "165545c9-2a5c-472c-8df2-7ff2be2b3b1b",
//   "Token" : "2336412f37fb687f5d51e6e241d09c805a5a57b30d712f794cc5f6a988666d92768dd60a747ba6f3beb71854e285d6ad02428b09ceece29417f1f02d609c582afbacc99c583a916b9981dd2728f4ae6fdb82efd087cc3b7849e05798d2d2785c03b0879594eeac82c01f235d0e717736",
//   "TopicArn" : "arn:aws:sns:us-west-2:123456789012:MyTopic",
//   "Message" : "You have chosen to subscribe to the topic arn:aws:sns:us-west-2:123456789012:MyTopic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
//   "SubscribeURL" : "https://sns.us-west-2.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-west-2:123456789012:MyTopic&Token=2336412f37fb687f5d51e6e241d09c805a5a57b30d712f794cc5f6a988666d92768dd60a747ba6f3beb71854e285d6ad02428b09ceece29417f1f02d609c582afbacc99c583a916b9981dd2728f4ae6fdb82efd087cc3b7849e05798d2d2785c03b0879594eeac82c01f235d0e717736",
//   "Timestamp" : "2012-04-26T20:45:04.751Z",
//   "SignatureVersion" : "1",
//   "Signature" : "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
//   "SigningCertURL" : "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem"
// }

// If the subscription succeeds we get a follow-up notification...
// // {
//   "Type" : "Notification",
//   "MessageId" : "228f8a83-4c07-5514-b316-dbeae8bbdb3f",
//   "TopicArn" : "arn:aws:sns:us-east-1:553053000801:accord_mojo_bounced_email",
//   "Message" : "{\"notificationType\":\"AmazonSnsSubscriptionSucceeded\",\"message\":\"You have successfully subscribed your Amazon SNS topic 'arn:aws:sns:us-east-1:553053000801:accord_mojo_bounced_email' to receive 'Bounce' notifications from Amazon SES for identity 'sman@accordinterests.com'.\"}\n",
//   "Timestamp" : "2017-04-13T16:08:28.706Z",
//   "SignatureVersion" : "1",
//   "Signature" : "dpA8h4FkzKCyKFIRNBvJMthI80NmYiN5jR8EzgQX88PZ7xjDT4PRHxb bUPHTCR0 GClXuNNhRUCUm XNod8zMgSura8cFw9mI7HAC8 rvV/6BZkfDcxI/ba1GSTXIGYC8dFjOp2IdIH0NBD7FjgT8VaBbp4rm9XJ3zTbkuC6 saybaY05I/Cv/xPrl/wFzlF urBhwmtEFLxq 8I5yjOpbmuZjLp9ejJRkQ6 5/tJ0hII16QVpY3OIpbLwRKSWOMJ2eZJll6XF1BdBUX3AWweh7KADHGpuQoNdZBJXrLKhsONilJZxHkV26nTboHC2yWAhl/46lHAJSjZ554UJbfQ==",
//   "SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem",
//   "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:553053000801:accord_mojo_bounced_email:4ed1756c-e31f-4809-a873-056c3026afb4"
// }

// AwsSubscribeConfirm is the struct of data sent by AWS
// to confirm a subscription to a topic.
type AwsSubscribeConfirm struct {
	Type             string    `json:"Type"`
	MessageID        string    `json:"MessageId"`
	Token            string    `json:"Token"`
	TopicArn         string    `json:"TopicArn"`
	Message          string    `json:"Message"`
	SubscribeURL     string    `json:"SubscribeURL"`
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
}

// AwsNotificationEnvelope describes the fields surrounding individual messages.
// We decode first into this struct to determin what type of message we received.
type AwsNotificationEnvelope struct {
	Type      string `json:"Type"`
	MessageID string `json:"MessageId"`
	TopicArn  string `json:"TopicArn"`
	Message   struct {
		NotificationType string `json:"notificationType"`
	}
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
	UnsubscribeURL   string    `json:"UnsubscribeURL"`
}

// SvcHandlerAws is the handler for aws subscription confirmation messages
func SvcHandlerAws(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcHandlerAws"
	fmt.Printf("Entered %s\n", funcname)

	// First thing to do is to look at the x-amz-sns-mssage-type header
	msgType := r.Header.Get("x-amz-sns-message-type")
	if len(msgType) == 0 {
		fmt.Printf("Could not find x-amz-sns-message-type header. Ignoring.\n")
		return
	}

	switch strings.ToLower(msgType) {
	case "subscriptionconfirmation":
		SvcHandlerAwsSubConf(w, r, d)
	case "notification":
		SvcHandlerNotification(w, r, d)
	default:
		fmt.Printf("Unhandled AWS message.  x-amz-sns-message-type = %s\n", msgType)
	}
}

// SvcHandlerNotification decodes and handles notifications from AWS
func SvcHandlerNotification(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcHandlerNotification"
	var a AwsNotificationEnvelope
	err := json.Unmarshal([]byte(d.data), &a)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		util.LogAndPrintError(funcname, e)
		return
	}
	fmt.Printf("Found Notification Type: %s\n", a.Message.NotificationType)
	switch a.Message.NotificationType {
	case "AmazonSnsSubscriptionSucceeded":
		util.Ulog("Notification Received: AmazonSnsSubscriptionSucceeded\n")
		fmt.Printf("Notification Received and processd: AmazonSnsSubscriptionSucceeded\n")
	case "Bounce":
		SvcHandlerAwsBouncedEmail(w, r, d)
	case "Complaint":
		SvcHandlerAwsComplaintEmail(w, r, d)
	default:
		fmt.Printf("Unhandled Notification Type: %s\n", a.Message.NotificationType)
	}
}

// SvcHandlerAwsSubConf is the handler for aws subscription confirmation messages
func SvcHandlerAwsSubConf(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcHandlerAwsSubConf"

	var a AwsSubscribeConfirm
	err := json.Unmarshal([]byte(d.data), &a)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		util.LogAndPrintError(funcname, e)
		return
	}
	// the proper reply is simply to respond with an HTTP Get on the URL in SubscribeURL
	fmt.Printf("HTTP GET to %s\n", a.SubscribeURL)
	response, err := http.Get(a.SubscribeURL)
	if err != nil {
		util.LogAndPrintError(funcname, err)
	} else {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			util.LogAndPrintError(funcname, err)
		}
		util.Ulog("Response from AWS to SubscribeURL: %s\n", string(body))
		fmt.Printf("Response: %s\n", string(body))
	}

}
