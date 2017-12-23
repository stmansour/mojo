package mailsend

import (
	"bytes"
	"fmt"
	"html/template"
	"mojo/db"
	"mojo/util"
	"rentroll/rlib"
	"strings"

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
	Hostname    string // host and domain of mojo server:   http://ex.domain.com:8275/
	SMTPHost    string // host of smtp server
	SMTPLogin   string // login name on smtp server
	SMTPPass    string // passwd for smtp login
	SMTPPort    int    // port to contact on smtp server
	Offset      int    // ignore if 0, otherwise use as query OFFSET
	Limit       int    // ignore if 0, otherwise use as query LIMIT
	DebugSend   bool   // if true, print the email addresses  but don't do the send
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

type pageData struct {
	P *db.Person
	L template.HTML // opt-out link
}

// GeneratePageHTML returns template.HTML for the supplied message, template, and person
//
// @params
//	fname    = name of the file with the base message html template
//  hostname = http(s)://hostname.domain:port/ where the opt-out link should point
//  p        = pointer to the db.Person record, the recipient of this message
//  t        = the html template
//
// @return
//	template.HTML = the body of the message
//  error         = any error that occured; nil on success
func GeneratePageHTML(fname, hostname string, p *db.Person, t *template.Template) (template.HTML, error) {
	funcname := "GeneratePageHTML"
	hostname = strings.TrimSuffix(hostname, "/") // remove last char if it is a slash.  it makes the Sprintf statement below easier to read.
	var pd pageData
	pd.P = p
	pd.L = template.HTML(fmt.Sprintf("%s/v1/optout?e=%s&c=%s", hostname, p.Email1, util.GenerateOptOutCode(p.FirstName, p.LastName, p.Email1, p.PID)))
	var sb bytes.Buffer
	err := t.Execute(&sb, &pd)
	if nil != err {
		util.LogAndPrintError(funcname, err)
		return template.HTML(""), err
	}
	return template.HTML(sb.String()), nil
}

// Sendmail is a routine to send an email message to a list of
// email addresses identified by the query.
func Sendmail(si *Info) error {
	funcname := "Sendmail"
	util.Ulog("MojoSendmail: Send message to %s\n", si.QName)
	util.Ulog("\t QName       = %s\n", si.QName)
	util.Ulog("\t From        = %s\n", si.From)
	util.Ulog("\t Subject     = %s\n", si.Subject)
	util.Ulog("\t MsgFName    = %s\n", si.MsgFName)
	util.Ulog("\t AttachFName = %s\n", si.AttachFName)
	util.Ulog("\t Hostname    = %s\n", si.Hostname)
	util.Ulog("\t SMTPHost    = %s\n", si.SMTPHost)
	util.Ulog("\t SMTPPort    = %d\n", si.SMTPPort)
	util.Ulog("\t OFFSET      = %d\n", si.Offset)
	util.Ulog("\t LIMIT       = %d\n", si.Limit)
	util.Ulog("\t DebugSend   = %t\n", si.DebugSend)

	// template for email
	var tname string
	sa := strings.Split(si.MsgFName, "/")
	if len(sa) > 0 {
		tname = sa[len(sa)-1]
	}

	t, err := template.New(tname).ParseFiles(si.MsgFName)
	if nil != err {
		fmt.Printf("%s: error loading template: %v\n", funcname, err)
		return err
	}

	// static portions of the message
	m := gomail.NewMessage()
	m.SetHeader("From", si.From)
	m.SetHeader("Subject", si.Subject)

	if len(si.AttachFName) > 0 {
		m.Attach(si.AttachFName)
	}

	q, err := BuildQuery(si.QName)
	if err != nil {
		e := fmt.Errorf("%s: Error getting database query: %s", funcname, err.Error())
		util.Ulog(e.Error() + "\n")
		return e
	}
	if si.Limit > 0 {
		q += fmt.Sprintf(" LIMIT %d", si.Limit)
	}
	if si.Offset > 0 {
		q += fmt.Sprintf(" OFFSET %d", si.Offset)
	}

	fmt.Printf("query is:\n%s\n", q)
	rows, err := db.DB.Db.Query(q)
	rlib.Errcheck(err)
	defer rows.Close()

	si.SentCount = 0
	good := 0      // number of valid email addresses
	bad := 0       // number of invalid email addresses
	optout := 0    // skipped because status was optout
	bounced := 0   // skipped because status was bounced
	complaint := 0 // skipped because status was complaint

	//fmt.Printf("EMAIL:  host: %s, port: %d, login: %s, pass: %s\n", si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
	d := gomail.NewDialer(si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
	for rows.Next() {
		p, err := db.ReadPersonFromRows(rows)
		if err != nil {
			util.Ulog("%s: Error with ReadPersonFromRows: %s\n", funcname, err.Error())
			return err
		}

		if len(p.Email1) == 0 {
			if si.DebugSend {
				fmt.Printf("%s: no email address for user: %d - %s %s\n", funcname, p.PID, p.FirstName, p.LastName)
			}
			continue
		}

		switch p.Status {
		case db.NORMAL:
			// do nothing
		case db.OPTOUT:
			optout++
			util.Ulog("%s: Email not sent to user %d, %s, due to prior: OPT OUT\n", funcname, p.PID, p.Email1)
			continue
		case db.BOUNCED:
			bounced++
			util.Ulog("%s: Email not sent to user %d, %s, due to prior: BOUNCED MESSAGE\n", funcname, p.PID, p.Email1)
			continue
		case db.COMPLAINT:
			complaint++
			util.Ulog("%s: Email not sent to user %d, %s, due to prior: COMPLAINT\n", funcname, p.PID, p.Email1)
			continue
		case db.SUPPRESSED:
			util.Ulog("%s: Email not sent to user %d, %s, due to prior: SUPPRESSION\n", funcname, p.PID, p.Email1)
			continue
		}

		//-------------------------------------------------------------------
		// We will always validate the email address first. They are often
		// completely bogus...
		//-------------------------------------------------------------------
		if !util.ValidEmailAddress(p.Email1) {
			bad++
			if si.DebugSend {
				fmt.Printf("%s: invalid email address %s for user: %d - %s %s\n", funcname, p.Email1, p.PID, p.FirstName, p.LastName)
			}
			util.Ulog("%s: Invalid email address:  PID = %d, email = %q\n", funcname, p.PID, p.Email1)
			continue
		} else {
			good++
		}

		m.SetHeader("To", p.Email1)
		s, err := GeneratePageHTML(si.MsgFName, si.Hostname, &p, t)
		if err != nil {
			util.Ulog("%s: Error on person with address: %s\n", funcname, p.Email1)
			return err
		}
		m.SetBody("text/html", string(s))

		if si.DebugSend {
			fmt.Printf("%s: Send to %s\n", funcname, p.Email1)
		} else {
			err = d.DialAndSend(m)
			if err != nil {
				util.Ulog("%s: Error on DialAndSend = %s\n", funcname, err.Error())
				util.Ulog("%s: Error occurred while sending to person %d, address: %s\n", funcname, p.PID, p.Email1)
				return err
			}
		}
		si.SentCount++ // update the si.SentCount only after adding the record

		if si.SentCount%25 == 0 {
			util.Console("%s: Processing query %s, SentCount = %d\n", funcname, si.QName, si.SentCount)
		}
	}

	util.Ulog("%s: Finished query %s\n", funcname, si.QName)
	util.Ulog("%s: Messages successfully sent:                 %5d\n", funcname, si.SentCount)
	util.Ulog("%s: Users skipped due to invalid email address: %5d\n", funcname, bad)
	util.Ulog("%s: Users skipped due to prior opt out:         %5d\n", funcname, optout)
	util.Ulog("%s: Users skipped due to prior bounced message: %5d\n", funcname, bounced)
	util.Ulog("%s: Users skipped due to prior complaint:       %5d\n", funcname, complaint)
	util.Ulog("%s: TOTAL USERS PROCESSED.......................%5d\n", funcname, si.SentCount+bad+optout+bounced+complaint)
	return nil
}
