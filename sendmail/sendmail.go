package sendmail

import (
	"fmt"
	"io/ioutil"
	"mojo/db"
	"mojo/util"
	"rentroll/rlib"

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
	SentCount   int    // count of messges sent
}

// BuildQuery creates a sql query from the JSON data
// stored in the supplied queryname
func BuildQuery(queryname string) (string, error) {
	q, err := db.GetQueryByName(queryname)
	if err != nil {
		return "", err
	}
	// TBD - call the actual translation code that xlates from
	// this structure into a sql query.  In the interim, I am
	// storing the actual sql query in the QueryJSON field
	// s := BuildSQLQuery(q.QueryJSON)
	return q.QueryJSON, nil // this is TEMPORARY
}

// Sendmail is a routine to send an email message to a list of
// email addresses identified by the query.
func Sendmail(si *Info) error {
	m := gomail.NewMessage()
	m.SetHeader("From", si.From)
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

	q, err := BuildQuery(si.QName)
	if err != nil {
		e := fmt.Errorf("Error getting database query: %s", err.Error())
		util.Ulog(e.Error() + "\n")
		return e
	}

	// fmt.Printf("query is %s\n", q)
	rows, err := db.DB.Db.Query(q)
	rlib.Errcheck(err)
	defer rows.Close()

	si.SentCount = 0
	for rows.Next() {
		p, err := db.ReadPersonFromRows(rows)
		if err != nil {
			util.Ulog("Error with ReadPersonFromRows: %s\n", err.Error())
			return err
		}
		m.SetHeader("To", p.Email1)
		// fmt.Printf("Sending to %s\n", p.Email1)
		d := gomail.NewDialer("email-smtp.us-east-1.amazonaws.com", 587, "AKIAJ3PENIYLS5U5ATJA", "AqIWufI4PwuxA61NihNQ4Yt+23n6w0CuQLuiUAdHP2E7")
		err = d.DialAndSend(m)
		if err != nil {
			util.Ulog("Error on DialAndSend = %s\n", err.Error())
			return err
		}
		si.SentCount++ // update the si.SentCount only after adding the record
	}
	return nil
}
