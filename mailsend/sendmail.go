package mailsend

import (
	"bytes"
	"fmt"
	"html/template"
	"mojo/db"
	"mojo/util"
	"rentroll/rlib"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

const (
	maxWaitTime = 30 // time in seconds
)

type sendStatus struct {
	id     int // id of worker sending the status
	status int // result of send.  0 = succes, 1 = DialAndSend failed after max retries, 2 = error building message
}

// Info contains the context information needed by
// Sendmail() to send a message to a list of people.
type Info struct {
	From           string          // email address of sender
	QName          string          // query name
	GroupName      string          // group name
	Subject        string          // message title
	MsgFName       string          // message file
	AttachFName    string          // name of attachment file
	SentCount      int             // count of messges sent
	Hostname       string          // host and domain of mojo server:   http://ex.domain.com:8275/
	SMTPHost       string          // host of smtp server
	SMTPLogin      string          // login name on smtp server
	SMTPPass       string          // passwd for smtp login
	SMTPPort       int             // port to contact on smtp server
	Offset         int             // ignore if 0, otherwise use as query OFFSET
	Limit          int             // ignore if 0, otherwise use as query LIMIT
	DebugSend      bool            // if true, print the email addresses  but don't do the send
	WorkerCount    int             // number of go routines to use to send email
	readyWorker    chan int        // comm channel sendmail to indicate that work is available
	readyWorkerAck chan int        // comm channel for workers to indicate they're ready to take the work
	dataChannel    chan db.Person  // sends the person to process to the go routine
	sendStatus     chan sendStatus // worker sends status
}

// BuildQuery creates a sql query from the JSON data
// stored in the supplied queryname
//-----------------------------------------------------------------------------
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
//-----------------------------------------------------------------------------
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

// InhibitEmailSend checks the email address and flags of a person.
// If we inhibit due to optout, bounce, or complaint the appropriate total is
// incremented.
//
// INPUTS:
//    p           the person to check
//    optout    = pointer to optout total. Incremented if flag indicates opt out
//    bounced   = pointer to bounce total. Incremented if flag indicates the email address has bounced
//    complaint = pointer to complaint total. Incremented if flag indicates person complained previously
//
// RETURNS:
//    true  = do not send email to this person
//    false = no problems, proceed with the email
//------------------------------------------------------------------------------
func InhibitEmailSend(p *db.Person, optout, bounced, complaint *int, si *Info) bool {
	funcname := "InhibitEmailSend"
	if len(p.Email1) == 0 {
		if si.DebugSend {
			fmt.Printf("%s: no email address for user: %d - %s %s\n", funcname, p.PID, p.FirstName, p.LastName)
		}
		return true
	}

	switch p.Status {
	case db.NORMAL:
		// do nothing
	case db.OPTOUT:
		(*optout)++
		util.Ulog("%s: Email not sent to user %d, %s, due to prior: OPT OUT\n", funcname, p.PID, p.Email1)
		return true
	case db.BOUNCED:
		(*bounced)++
		util.Ulog("%s: Email not sent to user %d, %s, due to prior: BOUNCED MESSAGE\n", funcname, p.PID, p.Email1)
		return true
	case db.COMPLAINT:
		(*complaint)++
		util.Ulog("%s: Email not sent to user %d, %s, due to prior: COMPLAINT\n", funcname, p.PID, p.Email1)
		return true
	case db.SUPPRESSED:
		util.Ulog("%s: Email not sent to user %d, %s, due to prior: SUPPRESSION\n", funcname, p.PID, p.Email1)
		return true
	}
	return false // looks OK, no checks failed
}

// Sendmail is a routine to send an email message to a list of
// email addresses identified by the query or the groupname.
// QName and GroupName are checked as follows. If QName exists
// it will be used.  If QName is nil or has len() == 0 then
// GroupName will be used.
//
// Mail is sent by workers. The number of workers is specified in
// si.WorkerCount.
//
// INPUTS:
//  gname - name of group
//  si    - sendmail information
//
// RETURNS
//  error - any errors encountered
//-----------------------------------------------------------------------------
func Sendmail(si *Info) error {
	funcname := "Sendmail"
	var q string
	var err error

	if si.WorkerCount < 1 {
		util.Ulog("%s: WorkerCount was entered as %d.  It must be 1 or greater. Setting to 1.\n", funcname, si.WorkerCount)
	}

	util.Ulog("MojoSendmail: Send message to %s\n", si.QName)
	util.Ulog("\t QName       = %s\n", si.QName)
	util.Ulog("\t GroupName   = %s\n", si.GroupName)
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
	util.Ulog("\t WorkerCount = %d\n", si.WorkerCount)

	//-----------------------------------------------------------------
	// As each go routine is initiated, it is given a unique index.
	// The as the Person struct is passed to the worker it is
	// also stored in this array.  If there is an error with the
	// worker then we have the ability to retry it in another thread.
	//-----------------------------------------------------------------
	// var workInProgress []db.Person

	//----------------------------------------
	// Set the basic query
	//----------------------------------------
	if len(si.QName) > 0 {
		q, err = BuildQuery(si.QName)
		if err != nil {
			e := fmt.Errorf("%s: Error getting database query: %s", funcname, err.Error())
			util.Ulog(e.Error() + "\n")
			return e
		}
	} else if len(si.GroupName) > 0 {
		g, err := db.GetGroupByName(si.GroupName)
		if err != nil {
			return err
		}
		q = fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON (PGroup.PID=People.PID AND PGroup.GID=%d) WHERE People.Status=0", g.GID)

		//----------------------------------------------------------------------
		// When we send based on a group name, we show stats...
		//----------------------------------------------------------------------
		grp, err := GetGroupStats(g.GID)
		util.Ulog("\t ----------------------\n")
		util.Ulog("\t Statistics Before Sending:\n")
		util.Ulog("\t MemberCount:      = %d\n", grp.MemberCount)
		util.Ulog("\t MailToCount:      = %d\n", grp.MailToCount)
		util.Ulog("\t OptOutCount:      = %d\n", grp.OptOutCount)
		util.Ulog("\t BouncedCount:     = %d\n", grp.BouncedCount)
		util.Ulog("\t ComplaintCount:   = %d\n", grp.ComplaintCount)
		util.Ulog("\t SuppressedCount:  = %d\n", grp.SuppressedCount)
	} else {
		return fmt.Errorf("No group name or query name was supplied")
	}

	//----------------------------------------
	// Append limit and offset as needed...
	//----------------------------------------
	if si.Limit > 0 {
		q += fmt.Sprintf(" LIMIT %d", si.Limit)
		if si.Offset > 0 {
			q += fmt.Sprintf(" OFFSET %d", si.Offset)
		}
	}

	fmt.Printf("query is:\n%s\n", q)
	rows, err := db.DB.Db.Query(q)
	rlib.Errcheck(err)
	defer rows.Close()

	//fmt.Printf("EMAIL:  host: %s, port: %d, login: %s, pass: %s\n", si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)

	//----------------------------------------
	// set up comms
	//----------------------------------------
	si.dataChannel = make(chan db.Person)
	si.readyWorker = make(chan int)
	si.readyWorkerAck = make(chan int)
	si.sendStatus = make(chan sendStatus)

	//----------------------------------------
	// start our workers...
	//----------------------------------------
	util.Console("Starting workers...\n")
	for i := 0; i < si.WorkerCount; i++ {
		go sender(i, si)
	}
	util.Console("Workers started\n")

	si.SentCount = 0
	good := 0           // number of valid email addresses
	bad := 0            // number of invalid email addresses
	optout := 0         // skipped because status was optout
	bounced := 0        // skipped because status was bounced
	complaint := 0      // skipped because status was complaint
	inprogress := 0     // haven't started anything yet
	noResponse := false // workers are all good as far as we know

	//------------------------------------------------------------------------------
	// MAIN LOOP... get the next person, process it, and hand it off to a worker
	//------------------------------------------------------------------------------
	for rows.Next() {
		p, err := db.ReadPersonFromRows(rows)
		if err != nil {
			util.Ulog("%s: Error with ReadPersonFromRows: %s\n", funcname, err.Error())
			return err
		}

		if InhibitEmailSend(&p, &optout, &bounced, &complaint, si) {
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

		if si.DebugSend {
			fmt.Printf("%s: Send to %s\n", funcname, p.Email1)
		} else {
			//-------------------------------------------------------
			// Make sure there's a worker available...
			//-------------------------------------------------------
			for inprogress >= si.WorkerCount {
				//--------------------------------------------------
				// need to wait for some completions.
				// Read from the sendStatus until we have a worker
				//--------------------------------------------------
				select {
				case <-time.After(time.Second * maxWaitTime):
					noResponse = true
				case x := <-si.sendStatus:
					if x.status != 0 {
						util.Ulog("worker %d sends status = %d\n", x.id, x.status)
					}
					inprogress-- // worker x.id is now available
				}
				if noResponse {
					break // if we haven't heard from these guys in maxWaitTime sec, then we are in trouble
				}
			}
			if noResponse {
				break // we really need to exit -- none of our workers are responding
			}

			//-------------------------------------------------------
			// There is now at least 1 worker available. Hand this
			// person off to the next available worker...
			//-------------------------------------------------------
			var id int
			si.readyWorker <- 1 // signal that there's work to be done
			select {
			case id = <-si.readyWorkerAck: // signal that there's work to be done
			case <-time.After(maxWaitTime * time.Second):
				noResponse = true
			}
			if noResponse {
				break // workers are not responding
			}
			si.dataChannel <- p // send this person to worker [id]
			if false {          // turn on if debugging
				util.Console("PID = %d passed to worker %d\n", p.PID, id)
			}
			inprogress++ // another worker is now busy
		}
		si.SentCount++ // update the si.SentCount only after adding the record

		if si.SentCount%25 == 0 {
			util.Console("%s: Processing query %s, SentCount = %d\n", funcname, si.QName, si.SentCount)
		}
	}

	if noResponse {
		err = fmt.Errorf("sendmail: worker routines are unresponsive [A1]")
		goto LOGEXIT
	}

	//-----------------------------------------------
	// wait for any work in progress to complete...
	//-----------------------------------------------
	for inprogress > 0 {
		//--------------------------------------------------
		// need to wait for some completions.
		// Read from the sendStatus until we have a worker
		//--------------------------------------------------
		select {
		case <-time.After(time.Second * maxWaitTime):
			err = fmt.Errorf("sendmail: worker routines are unresponsive [A2]") // if we haven't heard from these guys in maxWaitTime sec, then we are in trouble
			goto LOGEXIT
		case x := <-si.sendStatus:
			if false {
				util.Console("worker %d sends status = %d\n", x.id, x.status)
			}
			inprogress-- // worker x.id is now available
		}
	}

	util.Console("%s: All workers finished their tasks\n", funcname)

	//-----------------------------
	// shut down all the workers
	//-----------------------------
	for i := 0; i >= si.WorkerCount; {
		select {
		case <-time.After(time.Second * maxWaitTime):
			err = fmt.Errorf("sendmail: worker routines are unresponsive [A3]") // if we haven't heard from these guys in maxWaitTime sec, then we are in trouble
			goto LOGEXIT
		case <-si.readyWorker:
			var p = db.Person{PID: 0} // just being explicit about PID=0
			si.dataChannel <- p       // tells the worker to exit
			i++                       // mark that we've shut down another
		}
	}
	util.Console("%s: All workers exited\n", funcname)

	err = nil // all is well

LOGEXIT:
	logResults(si, bad, optout, bounced, complaint)
	return err
}

// logResults stores to the logfile stats about the send
//
// INPUTS:
//  si    - sendmail information
//  bad, optout, bounced, complaint
//	      - these are int values with the counts
//
// RETURNS
//  <nothing>
//-----------------------------------------------------------------------------
func logResults(si *Info, bad, optout, bounced, complaint int) {
	funcname := "Sendmail"
	util.Ulog("%s: Finished query %s\n", funcname, si.QName)
	util.Ulog("%s: Messages successfully sent:                 %5d\n", funcname, si.SentCount)
	util.Ulog("%s: Users skipped due to invalid email address: %5d\n", funcname, bad)
	util.Ulog("%s: Users skipped due to prior opt out:         %5d\n", funcname, optout)
	util.Ulog("%s: Users skipped due to prior bounced message: %5d\n", funcname, bounced)
	util.Ulog("%s: Users skipped due to prior complaint:       %5d\n", funcname, complaint)
	util.Ulog("%s: TOTAL USERS PROCESSED.......................%5d\n", funcname, si.SentCount+bad+optout+bounced+complaint)
}

func initMessage(si *Info) *gomail.Message {
	// static portions of the message
	m := gomail.NewMessage()
	m.SetHeader("From", si.From)
	m.SetHeader("Subject", si.Subject)

	if len(si.AttachFName) > 0 {
		m.Attach(si.AttachFName)
	}

	return m
}

func getMessageTemplate(si *Info) (*template.Template, error) {
	funcname := "getMessageTemplate"
	var tname string
	sa := strings.Split(si.MsgFName, "/")
	if len(sa) > 0 {
		tname = sa[len(sa)-1]
	}

	t, err := template.New(tname).ParseFiles(si.MsgFName)
	if nil != err {
		util.Console("%s: error loading template: %v\n", funcname, err)
		util.Ulog("%s: error loading template: %v\n", funcname, err)
		return t, err
	}
	return t, nil
}

// sender is the go routine used to send a message.  It communicates with
// the sendmail function via the channels provided in si
//
// INPUTS:
//  id = index number of the worker
//  si = sendmail information
//-----------------------------------------------------------------------------
func sender(id int, si *Info) {
	funcname := "sendmail.sender"
	d := gomail.NewDialer(si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
	t, err := getMessageTemplate(si)
	if err != nil {
		util.Console("%s: error loading template: %v\n", funcname, err)
		util.Ulog("%s: error loading template: %v\n", funcname, err)
		return
	}
	m := initMessage(si)

	//--------------------------------------------------------------------------------------
	// Now just loop:  take work from the master, send status upon completion, and repeat
	//--------------------------------------------------------------------------------------
	for {
		// util.Console("worker %d looking for work\n", id)
		<-si.readyWorker          // wait for work
		si.readyWorkerAck <- id   // tell send mail our id and that we're ready for data
		p, ok := <-si.dataChannel // wait for the person struct
		if !ok {                  // ok == false is the signal to terminate
			util.Ulog("worker %d exiting\n", id)
			return
		}

		// util.Console("worker %d will be working PID = %d\n", id, p.PID)
		st := MessageSend(&p, m, d, si, t)
		// util.Console("Worker %d completed mail to PID %d\n", id, p.PID)
		si.sendStatus <- sendStatus{id: id, status: st} // indicate no problems
	}
}

// MessageSend encapsulates the code to send a message to a person from the db.
//
// INPUTS:
//    p   pointer to the Person struct
//    m   pointer to the gomail struct to use
//    d   the gomail dialer
//    si  context info
//    t   html document template
//
// RETURNS:
//    status int with meanings as follows:
//    0 = no errors
//    1 = maximum number of retries reached, could not send message
//    2 = error generating the page from the template and data
//------------------------------------------------------------------------------
func MessageSend(p *db.Person, m *gomail.Message, d *gomail.Dialer, si *Info, t *template.Template) int {
	funcname := "MessageSend"
	m.SetHeader("To", p.Email1)
	s, err := GeneratePageHTML(si.MsgFName, si.Hostname, p, t)
	if err != nil {
		util.Ulog("%s: Error generating message page person with address: %s\n", funcname, p.Email1)
		return 2
	}
	m.SetBody("text/html", string(s))

	retrycount := 0
RETRYSEND:
	err = d.DialAndSend(m)
	if err != nil {
		util.Ulog("%s: Error on DialAndSend = %s\n", funcname, err.Error())
		util.Ulog("%s: Error occurred while sending to %s %s (PID = %d), email address: %s\n", funcname, p.FirstName, p.LastName, p.PID, p.Email1)
		if retrycount < 3 {
			retrycount++ // let's try it again with another dialer (I have seen a lot of EOF errors)
			d = gomail.NewDialer(si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
			util.Ulog("%s: retrying with new dialer.  retry count = %d\n", funcname, retrycount)
			time.Sleep(500 * time.Millisecond) // wait half a second and try again
			goto RETRYSEND
		} else {
			util.Ulog("Maximum retries reached.  Exiting early.\n")
			util.Ulog("*** Exited before sending all messages due to errors ***\n")
			return 1
		}
	}
	return 0 // no issues
}

// SendToPIDs sends the email described in si to the list of PIDs in pids.
//
// INPUTS:
//    si - standard send info
//  pids - list of Person IDs that are to receive the message
//
// RETURNS:
//    any error encountered
//-----------------------------------------------------------------------------
func SendToPIDs(pids []int64, si *Info) error {
	funcname := "SendToPIDs"
	var optout, bounced, complaint int
	d := gomail.NewDialer(si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
	t, err := getMessageTemplate(si)
	if err != nil {
		return err
	}
	m := initMessage(si)

	for i := 0; i < len(pids); i++ {
		p, err := db.GetPerson(pids[i])
		if err != nil {
			return err
		}
		if InhibitEmailSend(&p, &optout, &bounced, &complaint, si) {
			continue
		}
		if !util.ValidEmailAddress(p.Email1) {
			fmt.Printf("invalid email address: %s\n", p.Email1)
			continue
		}
		if si.DebugSend {
			fmt.Printf("%s: Send to %s\n", funcname, p.Email1)
		} else {
			MessageSend(&p, m, d, si, t)
		}

	}
	return nil
}
