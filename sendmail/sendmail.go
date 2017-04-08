package sendmail

import (
	"io/ioutil"
	"mojo/util"

	"gopkg.in/gomail.v2"
)

// Info contains the context information needed by
// Sendmail() to send a message to a list of people.
type Info struct {
	From        string // email address of sender
	QName       string // query name
	Subject     string // message title
	MsgFName    string // message file
	AttachFName string // name of attachment file
}

// Sendmail is a routine to send an email message to a list of
// email addresses identified by the query.
func Sendmail(si *Info) error {
	m := gomail.NewMessage()
	m.SetHeader("From", si.From)
	m.SetHeader("To", "sman@stevemansour.com")
	m.SetHeader("Subject", si.Subject)
	fileBytes, err := ioutil.ReadFile(si.MsgFName)
	if err != nil {
		util.Ulog("Error reading %s: %s\n", si.MsgFName, err.Error())
		return err
	}
	m.SetBody("text/html", string(fileBytes))
	if len(si.AttachFName) > 0 {
		m.Attach(si.AttachFName)
	}
	d := gomail.NewDialer("email-smtp.us-east-1.amazonaws.com", 587, "AKIAJ3PENIYLS5U5ATJA", "AqIWufI4PwuxA61NihNQ4Yt+23n6w0CuQLuiUAdHP2E7")
	return d.DialAndSend(m)
}
