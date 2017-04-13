package main

import (
	"encoding/json"
	"fmt"
	"mojo/util"
	"net/http"
	"time"
)

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

// SvcHandlerAwsComplaintEmail removes a bounced email address from the database
func SvcHandlerAwsComplaintEmail(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcHandlerAwsComplaintEmail"
	fmt.Printf("Entered %s\n", funcname)
	var a AwsComplaintNotification
	err := json.Unmarshal([]byte(d.data), &a)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		util.LogAndPrintError(funcname, e)
		return
	}

	fmt.Printf("Received Complaint Email Message!\n")
	fmt.Printf("%#v\n", a)

}
