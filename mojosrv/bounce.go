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
	NotificationType string `json:"notificationType"`
	Bounce           struct {
		BounceType        string `json:"bounceType"`
		ReportingMTA      string `json:"reportingMTA"`
		BouncedRecipients []struct {
			EmailAddress   string `json:"emailAddress"`
			Status         string `json:"status"`
			Action         string `json:"action"`
			DiagnosticCode string `json:"diagnosticCode"`
		} `json:"bouncedRecipients"`
		BounceSubType string    `json:"bounceSubType"`
		Timestamp     time.Time `json:"timestamp"`
		FeedbackID    string    `json:"feedbackId"`
		RemoteMtaIP   string    `json:"remoteMtaIp"`
	} `json:"bounce"`
	Mail AwsMailNotification `json:"mail"`
}

// AwsComplaintNotification is the data type for an AWS complaint email message notification
type AwsComplaintNotification struct {
	NotificationType string `json:"notificationType"`
	Complaint        struct {
		UserAgent            string `json:"userAgent"`
		ComplainedRecipients []struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"complainedRecipients"`
		ComplaintFeedbackType string    `json:"complaintFeedbackType"`
		ArrivalDate           time.Time `json:"arrivalDate"`
		Timestamp             time.Time `json:"timestamp"`
		FeedbackID            string    `json:"feedbackId"`
	} `json:"complaint"`
	Mail AwsMailNotification `json:"mail"`
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
	fmt.Printf("%#v\n", a)

	for i := 0; i < len(a.Mail.CommonHeaders.To); i++ {
		fmt.Printf("Email address to remove: %s\n", a.Mail.CommonHeaders.To[i])
	}
}
