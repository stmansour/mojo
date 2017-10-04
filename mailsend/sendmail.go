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
		e := fmt.Errorf("Error getting database query: %s", err.Error())
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
	//fmt.Printf("EMAIL:  host: %s, port: %d, login: %s, pass: %s\n", si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
	d := gomail.NewDialer(si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
	for rows.Next() {
		p, err := db.ReadPersonFromRows(rows)
		if err != nil {
			util.Ulog("Error with ReadPersonFromRows: %s\n", err.Error())
			return err
		}
		// fmt.Printf("Sending to %s\n", p.Email1)
		if len(p.Email1) == 0 {
			continue
		}
		m.SetHeader("To", p.Email1)
		s, err := GeneratePageHTML(si.MsgFName, si.Hostname, &p, t)
		if err != nil {
			util.Ulog("Error on person with address: %s\n", p.Email1)
			return err
		}
		m.SetBody("text/html", string(s))

		err = d.DialAndSend(m)
		if err != nil {
			util.Ulog("Error on DialAndSend = %s\n", err.Error())
			util.Ulog("Error occurred while sending to person %d, address: %s\n", p.PID, p.Email1)
			return err
		}
		si.SentCount++ // update the si.SentCount only after adding the record

		if si.SentCount%25 == 0 {
			util.Ulog("Processing query %s, SentCount = %d\n", si.QName, si.SentCount)
		}
	}

	util.Ulog("Finished query %s. Successfully sent %d messages\n", si.QName, si.SentCount)
	return nil
}
