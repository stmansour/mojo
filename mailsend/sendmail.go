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

func generatePageHTML(fname, hostname string, p *db.Person, t *template.Template) (template.HTML, error) {
	funcname := "generatePageHTML"
	hostname = strings.TrimSuffix(hostname, "/") // remove last char if it is a slash.  it makes the Sprintf statement below easier to read.
	var pd pageData
	pd.P = p
	pd.L = template.HTML(fmt.Sprintf("%s/v1/optout?e=%s&c=%s", hostname, p.Email1, util.GenerateOptOutCode(p.FirstName, p.LastName, p.Email1, p.PID)))
	var sb bytes.Buffer
	err := t.Execute(&sb, &pd)
	if nil != err {
		util.LogAndPrintError(funcname, err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	//s := sb.String()
	return template.HTML(sb.String()), nil
}

// Sendmail is a routine to send an email message to a list of
// email addresses identified by the query.
func Sendmail(si *Info) error {
	funcname := "Sendmail"
	util.Ulog("MojoSendmail: Calling %s to send message to List - query = %s, from = %s\n", si.Hostname, si.QName, si.From)

	// template for email
	t, err := template.New(si.MsgFName).ParseFiles(si.MsgFName)
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
		s, err := generatePageHTML(si.MsgFName, si.Hostname, &p, t)
		if err != nil {
			return err
		}
		m.SetBody("text/html", string(s))

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
