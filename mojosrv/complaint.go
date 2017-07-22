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
func SvcHandlerAwsComplaintEmail(w http.ResponseWriter, r *http.Request, d *ServiceData, a *AwsNotificationEnvelope) {
	funcname := "SvcHandlerAwsComplaintEmail"
	util.Console("Entered %s\n", funcname)
	var b AwsComplaintNotification
	err := json.Unmarshal([]byte(a.Message), &b)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		util.LogAndPrintError(funcname, e)
		return
	}
	util.Console("Complaint recipient email addresses:\n")
	for i := 0; i < len(b.Complaint.ComplainedRecipients); i++ {
		util.Console("%d. %s\n", i, b.Complaint.ComplainedRecipients[i].EmailAddress)
		err = HandleEmailComplaint(b.Complaint.ComplainedRecipients[i].EmailAddress)
		if err != nil {
			util.Ulog("%s: Error handling email address %s: %s\n",
				funcname, b.Complaint.ComplainedRecipients[i].EmailAddress, err.Error())
		}
	}
}
