package main

import (
	"encoding/json"
	"fmt"
	"mojo/db"
	"mojo/util"
	"net/http"
	"time"
)

// This module handles the bounce messages from AWS SNS.
//
// For reference:
// 	bounce object: http://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html#bounce-object
// 	complaint object: http://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html#complaint-object
// 	delivery object: http://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html#delivery-object

// AwsMailNotification is the data type associated with the AWS mail notification
type AwsMailNotification struct {
	Timestamp        time.Time `json:"timestamp"`
	Source           string    `json:"source"`
	SourceArn        string    `json:"sourceArn"`
	SourceIP         string    `json:"sourceIp"`
	SendingAccountID string    `json:"sendingAccountId"`
	MessageID        string    `json:"messageId"`
	Destination      []string  `json:"destination"`
	HeadersTruncated bool      `json:"headersTruncated"`
	Headers          []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"headers"`
	CommonHeaders struct {
		From      []string `json:"from"`
		Date      string   `json:"date"`
		To        []string `json:"to"`
		MessageID string   `json:"messageId"`
		Subject   string   `json:"subject"`
	} `json:"commonHeaders"`
}

// AwsBounceNotification is the data type for an AWS Bounced email message notification
type AwsBounceNotification struct {
	Type      string `json:"Type"`
	MessageID string `json:"MessageId"`
	TopicArn  string `json:"TopicArn"`
	Message   struct {
		NotificationType string `json:"notificationType"`
		Bounce           struct {
			BounceType        string `json:"bounceType"`
			BounceSubType     string `json:"bounceSubType"`
			BouncedRecipients []struct {
				EmailAddress   string `json:"emailAddress"`
				Action         string `json:"action"`
				Status         string `json:"status"`
				DiagnosticCode string `json:"diagnosticCode"`
			} `json:"bouncedRecipients"`
			Timestamp    time.Time `json:"timestamp"`
			FeedbackID   string    `json:"feedbackId"`
			RemoteMtaIP  string    `json:"remoteMtaIp"`
			ReportingMTA string    `json:"reportingMTA"`
		} `json:"bounce"`
		Mail struct {
			Timestamp        time.Time `json:"timestamp"`
			Source           string    `json:"source"`
			SourceArn        string    `json:"sourceArn"`
			SourceIP         string    `json:"sourceIp"`
			SendingAccountID string    `json:"sendingAccountId"`
			MessageID        string    `json:"messageId"`
			Destination      []string  `json:"destination"`
		} `json:"mail"`
	} `json:"Message"`
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
	UnsubscribeURL   string    `json:"UnsubscribeURL"`
}

// ChangePersonStatus is called with the email address of the person
// to be updated. If found, the person's status will be set to the
// supplied value.
func ChangePersonStatus(s string, status int64) error {
	p, err := db.GetPersonByEmail(s)
	if err != nil {
		if !util.IsSQLNoResultsError(err) {
			util.Ulog("ChangePersonStatus: error getting Person for %s")
		}
		return err
	}
	p.Status = status
	return db.UpdatePerson(&p)
}

// HandleEmailBounce is called with the bounced email address. It will
// update the associated person record with a status of BOUNCED
func HandleEmailBounce(s string) error {
	return ChangePersonStatus(s, db.BOUNCED)
}

// HandleEmailComplaint is called with the bounced email address. It will
// update the associated person record with a status of COMPLAINT
func HandleEmailComplaint(s string) error {
	return ChangePersonStatus(s, db.COMPLAINT)
}

// SvcHandlerAwsBouncedEmail removes a bounced email address from the database
func SvcHandlerAwsBouncedEmail(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcHandlerAwsBouncedEmail"
	fmt.Printf("Entered %s\n", funcname)
	var a AwsBounceNotification
	err := json.Unmarshal([]byte(d.data), &a)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		util.LogAndPrintError(funcname, e)
		return
	}

	fmt.Printf("Received Bounced Email Message!\n")
	// fmt.Printf("%#v\n", a)

	for i := 0; i < len(a.Message.Bounce.BouncedRecipients); i++ {
		fmt.Printf("Email address to remove: %s\n", a.Message.Bounce.BouncedRecipients[i].EmailAddress)
	}
}
