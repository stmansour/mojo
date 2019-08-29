package main

import (
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"log"
	"mojo/db"
	"mojo/mailsend"
	"mojo/util"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	_ "github.com/go-sql-driver/mysql"
)

// App is the application structure available to the whole app.
var App struct {
	MsgFile       string
	AttachFile    string
	QueryName     string
	GroupName     string
	MojoHost      string // domain and port for mojosrv:   http://example.domain.com:8275/
	db            *sql.DB
	DBName        string
	DBUser        string
	ValidateGroup string // validate all the email addresses in this group
	SetupOnly     bool
	Subject       string                   // subject line of the email message
	From          string                   // email from address
	PIDs          []int64                  // array of PIDs to send to
	QueryCount    bool                     // if true, just print the solution set count for the query and exit
	Bounce        bool                     // if true, just print the solution set count for the query and exit
	Complaint     bool                     // if true, just print the solution set count for the query and exit
	OOO           bool                     // if true, just print the solution set count for the query and exit
	Suppress      bool                     // if true, just print the solution set count for the query and exit
	Fix           bool                     // if true, just scan the database for errors, fix the ones we can, then exit
	LogFile       *os.File                 // where to log messages
	XR            extres.ExternalResources // dbs, smtp...
	Offset        int                      // if 0, ignore, if nonzero then use as query OFFSET
	Limit         int                      // if 0, ignore, if nonzero then use as query LIMIT
	DebugSend     bool                     // print email addresses where send would have gone but don't send
	WorkerCount   int                      // number of go routines to use to send the email; parallelization
}

// SendBouncedEmailTest sends an email message that bounces.  For testing.
// The recipient's ISP rejects your email with an SMTP 550 5.1.1 response
// code ("Unknown User"). Amazon SES generates a bounce notification and
// sends it to you via email or by using an Amazon SNS notification,
// depending on how you set up your system. This mailbox simulator email
// address will not be placed on the Amazon SES suppression list as one
// normally would when an email hard bounces. The bounce response that
// you receive from the mailbox simulator is compliant with RFC 3464.
//-----------------------------------------------------------------------------
func SendBouncedEmailTest() error {
	return SendEmailTest("bounce@simulator.amazonses.com")
}

// SendComplaintEmailTest sends an email message that bounces.  For testing.
// The recipient's ISP accepts your email and delivers it to the recipient’s
// inbox. The recipient, however, does not want to receive your message and
// clicks "Mark as Spam" within an email application that uses an ISP that
// sends a complaint response to Amazon SES. Amazon SES then forwards the
// complaint notification to you via email or by using an Amazon SNS
// notification, depending on how you set up your system. The complaint
// response that you receive from the mailbox simulator is compliant with
// RFC 5965.
//-----------------------------------------------------------------------------
func SendComplaintEmailTest() error {
	return SendEmailTest("complaint@simulator.amazonses.com")
}

// SendOOOEmailTest sends a test email. The recipient's ISP accepts your
// email and delivers it to the recipient’s inbox. The ISP sends an
// out-of-the-office (OOTO) message to Amazon SES. Amazon SES then forwards
// the OOTO message to you via email or by using an Amazon SNS notification,
// depending on how you set up your system. The OOTO response that you receive
// from the Mailbox Simulator is compliant with RFC 3834. For information
// about how to set up your system to receive OOTO responses, follow the same
// instructions for setting up how Amazon SES sends you notifications in
// Monitoring Using Amazon SES Notifications.
//-----------------------------------------------------------------------------
func SendOOOEmailTest() error {
	return SendEmailTest("ooto@simulator.amazonses.com")
}

// SendSuppressionListEmailTest sends a test email. Amazon SES treats your
// email as a hard bounce because the address you are sending to is on the
// Amazon SES suppression list.
//-----------------------------------------------------------------------------
func SendSuppressionListEmailTest() error {
	return SendEmailTest("suppressionlist@simulator.amazonses.com")
}

// SendEmailTest is a routine to send an email message to the supplied address.
//-----------------------------------------------------------------------------
func SendEmailTest(addr string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "sman@stevemansour.com")
	m.SetHeader("Subject", "Force a bounce message")
	m.SetBody("text/html", "<html><body><p>This should bounce!</p></body></html>")
	m.SetHeader("To", addr)
	fmt.Printf("Sending BOUNCE message to %s\n", addr)
	// d := gomail.NewDialer("email-smtp.us-east-1.amazonaws.com", 587, "AKIAJ3PENIYLS5U5ATJA", "AqIWufI4PwuxA61NihNQ4Yt+23n6w0CuQLuiUAdHP2E7")
	d := gomail.NewDialer(db.MojoDBConfig.SMTPHost, db.MojoDBConfig.SMTPPort, db.MojoDBConfig.SMTPLogin, db.MojoDBConfig.SMTPPass)

	err := d.DialAndSend(m)
	if err != nil {
		util.Ulog("Error on DialAndSend = %s\n", err.Error())
		return err
	}
	fmt.Printf("Bounce message successfully sent to %s\n", addr)
	return nil
}

// AddPersonToGroup creates a PGroup record for the specified pid,gid pair
// if it does not already exist.
//-----------------------------------------------------------------------------
func AddPersonToGroup(pid, gid int64) error {
	// see if they already exist...
	_, err := db.GetPGroup(pid, gid)
	if util.IsSQLNoResultsError(err) {
		var a = db.PGroup{PID: pid, GID: gid}
		err = db.InsertPGroup(&a)
		if err != nil {
			util.Ulog("Error with InsertPGroup: %s\n", err.Error())
		}
		return err
	}
	if err == nil {
		return nil // they're already in the group
	}
	util.Ulog("Error trying to GetPGroup = %s\n", err.Error())
	return err
}

// SavePerson creates a new person in the database with the supplied
// information. If the person already exists, it updates their info
// with whatever is in pnew.
//
// INPUTS
//  pnew = struct with information about a person to save
//
// RETURNS
//  pid of the person if error is nil
//  any error encountered
//-----------------------------------------------------------------------------
func SavePerson(pnew *db.Person) (int64, error) {
	var pid int64
	p1, err := db.GetPersonByName(pnew.FirstName, pnew.MiddleName, pnew.LastName)
	if err != nil {
		util.Ulog("db.GetPersonByName returned: %s\n", err.Error())
		return pid, err
	}
	if len(p1) == 0 {
		err := db.InsertPerson(pnew)
		if err != nil {
			util.Ulog("db.InsertPerson returned: %s\n", err.Error())
			return pid, err
		}
		pid = pnew.PID
	} else {
		pid = p1[0].PID
		p1[0].FirstName = pnew.FirstName
		p1[0].MiddleName = pnew.MiddleName
		p1[0].LastName = pnew.LastName
		p1[0].PreferredName = pnew.PreferredName
		p1[0].JobTitle = pnew.JobTitle
		p1[0].OfficePhone = pnew.OfficePhone
		p1[0].OfficeFax = pnew.OfficeFax
		p1[0].Email1 = pnew.Email1
		p1[0].Email2 = pnew.Email2
		p1[0].MailAddress = pnew.MailAddress
		p1[0].MailAddress2 = pnew.MailAddress2
		p1[0].MailCity = pnew.MailCity
		p1[0].MailState = pnew.MailState
		p1[0].MailPostalCode = pnew.MailPostalCode
		p1[0].MailCountry = pnew.MailCountry
		p1[0].RoomNumber = pnew.RoomNumber
		p1[0].MailStop = pnew.MailStop
		err = db.UpdatePerson(&p1[0])
		if err != nil {
			return pid, err
		}
	}
	return pid, nil
}

func resetUserStatus(p1 *db.Person) {
	// util.Console("Entered resetUserStatus - person with email: %s\n", p1.Email1)
	p, err := db.GetPersonByEmail(p1.Email1)
	if err != nil {
		util.UlogAndPrint("Error from db.GetPersonByEmail( %s ):  %s \n", p1.Email1, err.Error())
		return
	}
	// util.Console("Found person: %s %s   Status = %d\n", p.FirstName, p.LastName, p.Status)
	p.Status = db.NORMAL
	// util.Console("Updated status to %d\n", p.Status)
	if err = db.UpdatePerson(&p); err != nil {
		util.UlogAndPrint("Error from db.UpdatePerson( %s ):  %s \n", p.Email1, err.Error())
		return
	}
	// util.Console("Update successful\n")
}

func createGroup(name, descr string, ppa *[]db.Person) {
	var g db.EGroup

	g, err := db.GetGroupByName(name)
	if err == nil {
		fmt.Printf("%s exists, updating timestamp\n", name)
		g.DtStart = time.Now()
		g.DtStop = time.Now()
		err = db.UpdateGroup(&g)
		if err != nil {
			util.UlogAndPrint("Error updating group: %s\n", err.Error())
			os.Exit(1)
		}
	}
	if err != nil {
		if !util.IsSQLNoResultsError(err) {
			util.UlogAndPrint("Error reading group \"MojoTest\": %s\n", err.Error())
			os.Exit(1)
		}
	}

	// Create the group...
	if g.GID == 0 {
		g.GroupName = name
		g.GroupDescription = descr
		g.DtStart = time.Now()
		err = db.InsertGroup(&g)
		if err != nil {
			util.UlogAndPrint("Error inserting group: %s\n", err.Error())
			os.Exit(1)
		}
	}

	// Add the list of people to it...
	pa := *ppa
	gid := g.GID
	for i := 0; i < len(pa); i++ {
		pid, err := SavePerson(&pa[i]) // if person is not in db, add them, then add them to group gid
		if err != nil {
			util.UlogAndPrint("Error updating group: %s\n", err.Error())
		}
		AddPersonToGroup(pid, gid)
		// util.Console("Reset user status for %s %s\n", pa[i].FirstName, pa[i].LastName)
		resetUserStatus(&pa[i]) // reset their status
	}

	// Mark that we're finished...
	g.DtStop = time.Now()
	err = db.UpdateGroup(&g)
	if err != nil {
		util.UlogAndPrint("Error updating group: %s\n", err.Error())
		os.Exit(1)
	}

	// create a default query
	q := fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0", g.GID)
	createQuery(name, descr, q)
}

func createQuery(name, descr, query string) {
	q, err := db.GetQueryByName(name)
	if err != nil {
		if util.IsSQLNoResultsError(err) {
			q.QueryName = name
			q.QueryDescr = descr
			q.QueryJSON = query
			err = db.InsertQuery(&q)
			if err != nil {
				util.UlogAndPrint("Error inserting query: %s\n", err.Error())
				os.Exit(1)
			}
		} else {
			util.UlogAndPrint("Error reading query %q: %s\n", name, err.Error())
			os.Exit(1)
		}
	} else {
		if query != q.QueryJSON {
			q.QueryJSON = query
			err = db.UpdateQuery(&q)
			if err != nil {
				util.UlogAndPrint("Error updating query: %s\n", err.Error())
				os.Exit(1)
			}
		}
	}
}

func createIsolaBellaQueries() {
	var g db.EGroup
	var err error
	grp := "OK Real Estate Agents"
	g, err = db.GetGroupByName(grp)
	if err != nil && util.IsSQLNoResultsError(err) {
		fmt.Printf("Group %q does not exist... no queries added\n", grp)
		return
	}
	if err != nil {
		util.UlogAndPrint("Error getting group %q: %s\n", grp, err.Error())
		os.Exit(1)
	}
	q := fmt.Sprintf("SELECT People.* from People INNER JOIN PGroup ON (PGroup.PID=People.PID AND PGroup.GID=%d) WHERE People.Status=0 AND (People.MailPostalCode LIKE \"73750%%\" OR People.MailPostalCode LIKE \"73762%%\" OR People.MailPostalCode LIKE \"73078%%\" OR People.MailPostalCode LIKE \"73016%%\" OR People.MailPostalCode LIKE \"73735%%\" OR People.MailPostalCode LIKE \"73773%%\" OR People.MailPostalCode LIKE \"73718%%\" OR People.MailPostalCode LIKE \"73763%%\")", g.GID)
	createQuery("ISO-Frack", "includes zipcodes 73750...", q)
	q = fmt.Sprintf("SELECT People.* from People INNER JOIN PGroup ON (PGroup.PID=People.PID AND PGroup.GID=%d) WHERE People.Status=0 AND NOT (People.MailPostalCode LIKE \"73750%%\" OR People.MailPostalCode LIKE \"73762%%\" OR People.MailPostalCode LIKE \"73078%%\" OR People.MailPostalCode LIKE \"73016%%\" OR People.MailPostalCode LIKE \"73735%%\" OR People.MailPostalCode LIKE \"73773%%\" OR People.MailPostalCode LIKE \"73718%%\" OR People.MailPostalCode LIKE \"73763%%\")", g.GID)
	createQuery("ISO-Commissions", "excludes zipcodes 73750...", q)

	g, err = db.GetGroupByName("FAA Tech Ops")
	if err != nil && util.IsSQLNoResultsError(err) {
		fmt.Printf("Group %q does not exist... no queries added\n", grp)
		return
	}
	if err != nil {
		util.UlogAndPrint("Error getting group %q: %s\n", grp, err.Error())
		os.Exit(1)
	}
	q = fmt.Sprintf("SELECT People.* from People INNER JOIN PGroup ON (PGroup.PID=People.PID AND PGroup.GID=%d) WHERE People.Status=0", g.GID)
	createQuery("ISO-FAATechOps", "FAA employees who have stayed at IB", q)

	g, err = db.GetGroupByName("ibguests20171206")
	if err != nil && util.IsSQLNoResultsError(err) {
		fmt.Printf("Group %q does not exist... no queries added\n", grp)
		return
	}
	if err != nil {
		util.UlogAndPrint("Error getting group %q: %s\n", grp, err.Error())
		os.Exit(1)
	}
	q = fmt.Sprintf("SELECT People.* from People INNER JOIN PGroup ON (PGroup.PID=People.PID AND PGroup.GID=%d) WHERE People.Status=0", g.GID)
	createQuery("ibguests20171206", "FAA employees who have stayed at IB", q)
	fmt.Printf("Isola Bella queries created\n")
}

func createFAAQueries() {
	var g db.EGroup
	var err error

	g, err = db.GetGroupByName("FAA")
	if err != nil && util.IsSQLNoResultsError(err) {
		fmt.Printf("Group FAA does not exist... no FAA queries added\n")
		return
	}
	if err != nil {
		util.UlogAndPrint("Error getting group FAA: %s\n", err.Error())
		os.Exit(1)
	}
	q := fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0", g.GID) // 50
	createQuery("FAA", "The first 50 people in the FAA", q)
}

// Set up the people information first. This will make the people available
// for creating groups, for "setup" only (that is, reset status), and it will
// make it easy to add new people to the groups.
//-----------------------------------------------------------------------------
var pa = []db.Person{
	{FirstName: "Steven", MiddleName: "F", LastName: "Mansour", JobTitle: "CTO, Accord Interests", OfficePhone: "323-512-0111 X305", Email1: "sman@accordinterests.com", MailAddress: "11719 Bee Cave Road", MailAddress2: "Suite 301", MailCity: "Austin", MailState: "TX", MailPostalCode: "78738", MailCountry: "USA", Status: 0},
	{FirstName: "Steve", MiddleName: "", LastName: "Mansour", JobTitle: "Recording Musician, Engineer, Producer", OfficePhone: "323-512-0111 X305", Email1: "sman@stevemansour.com", MailAddress: "2215 Wellington Drive", MailAddress2: "", MailCity: "Milpitas", MailState: "CA", MailPostalCode: "95035", MailCountry: "USA", Status: 0},
}

var pa1 = []db.Person{
	pa[0],
	pa[1],
	{FirstName: "Bouncie", MiddleName: "", LastName: "McBounce", JobTitle: "Vagabond", OfficePhone: "123-456-7890", Email1: "bounce@simulator.amazonses.com", MailAddress: "123 Elm St", MailAddress2: "", MailCity: "Anytown", MailState: "CA", MailPostalCode: "90210", MailCountry: "USA", Status: 0},
	{FirstName: "Wendy", MiddleName: "", LastName: "Whiner", JobTitle: "Complainer", OfficePhone: "123-321-7890", Email1: "complaint@simulator.amazonses.com", MailAddress: "321 Elm St", MailAddress2: "", MailCity: "Anytown", MailState: "CA", MailPostalCode: "90210", MailCountry: "USA", Status: 0},
	{FirstName: "Stealthy", MiddleName: "", LastName: "McStealth", JobTitle: "Bad Guy", OfficePhone: "816-321-0123", Email1: "suppressionlist@simulator.amazonses.com", MailAddress: "700 Elm St", MailAddress2: "", MailCity: "Anytown", MailState: "CA", MailPostalCode: "90210", MailCountry: "USA", Status: 0},
}
var pa2 = []db.Person{
	pa1[0],
	pa1[1],
	pa1[2],
	pa1[3],
	pa1[4],
	{FirstName: "Joe", MiddleName: "G", LastName: "Mansour", JobTitle: "Principal, Accord Interests", OfficePhone: "323-512-0111 X303", Email1: "jgm@accordinterests.com", MailAddress: "11719 Bee Cave Road", MailAddress2: "Suite 301", MailCity: "Austin", MailState: "TX", MailPostalCode: "78738", MailCountry: "USA", Status: 0},
	{FirstName: "Melissa", MiddleName: "", LastName: "Wheeler", JobTitle: "General Manager, Isola Bella", OfficePhone: "405.721.2194 x205", Email1: "mwheeler@accordinterests.com.com", MailAddress: "8309 NW 140th St", MailAddress2: "", MailCity: "Oklahoma City", MailState: "OK", MailPostalCode: "73142", MailCountry: "USA", Status: 0},
	//	{FirstName: "Michelle", MiddleName: "", LastName: "Falls", JobTitle: "Concierge", OfficePhone: "405.721.2194 x2014", Email1: "mfalls@myisolabella.com", MailAddress: "8309 NW 140th St", MailAddress2: "", MailCity: "Oklahoma City", MailState: "OK", MailPostalCode: "73142", MailCountry: "USA", Status: 0},
	// {FirstName: "Brittney", MiddleName: "", LastName: "Graham", JobTitle: "Manager", OfficePhone: "405.721.2194", Email1: "bgraham@myisolabella.com", MailAddress: "6608 Lyrewood Ln", MailAddress2: "Apt 24", MailCity: "Oklahoma City", MailState: "OK", MailPostalCode: "73132", MailCountry: "USA", Status: 0},
	//	{FirstName: "Kristy", MiddleName: "", LastName: "Koon", JobTitle: "Serviced Apt Sales", OfficePhone: "405.721.2194", Email1: "kkoon@myisolabella.com", MailAddress: "10407 SE 23rd", MailAddress2: "Apt 24", MailCity: "Oklahoma City", MailState: "OK", MailPostalCode: "73130", MailCountry: "USA", Status: 0},
}

var pa3 = []db.Person{
	pa[0],
	pa[1],
	pa2[5],
}

func setupTestGroups() {
	//--------------------------------------------------------
	// Make sure that all people are in the database first...
	//--------------------------------------------------------
	for i := 0; i < len(pa2); i++ {
		_, err := SavePerson(&pa2[i]) // if person is not in db, add them, then add them to group gid
		if err != nil {
			util.LogAndPrintError("setupTestGroups", err)
			os.Exit(1)
		}
		resetUserStatus(&pa2[i]) // reset their status
	}

	createGroup("MojoTest", "Steve-only test group", &pa)
	createGroup("AmazonTest", "Steve + Amazon test accounts", &pa1)
	createGroup("AccordTest", "Steve + Amazon test + Accord accounts", &pa2)
	createGroup("SteveJoe", "Steve + Joe accounts", &pa3)
	createIsolaBellaQueries()
	createFAAQueries()
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "mojo", "database name")
	mPtr := flag.String("b", "testmsg.html", "filename containing the html message to send")
	aPtr := flag.String("a", "", "filename of attachment")
	qPtr := flag.String("q", "", "name of the query to send messages to; overrides group if supplied")
	hPtr := flag.String("h", db.MojoDBConfig.MojoWebAddr, "name of host and port for mojosrv")
	vPtr := flag.String("validate", "", "validate the email addresses of everyone in the group name provided, then exit")
	gPtr := flag.String("group", "", "group name to send mail to; overridden by -q if it is supplied")
	pids := flag.String("pids", "", "comma separated list of PIDs to send to, overrides -p and -q, ex: -pids 764,5263,3452")
	qcPtr := flag.Bool("count", false, "returns the count of target addresses in the query, then exits.")
	soPtr := flag.Bool("setup", false, "just run the setup, do not send email")
	bPTR := flag.Bool("bounce", false, "just send a message to bounce@simulator.amazonses.com")
	cPTR := flag.Bool("complaint", false, "just send a message to complaint@simulator.amazonses.com")
	oPTR := flag.Bool("ooo", false, "just send a message to ooo@simulator.amazonses.com")
	sPTR := flag.Bool("sl", false, "just send a message to suppressionlist@simulator.amazonses.com")
	subjPtr := flag.String("subject", "Test Message", "Email subject line.")
	fromPtr := flag.String("from", "sman@accordinterests.com", "Message sender.")
	fixPtr := flag.Bool("fix", false, "Scan db for known errors, fix them wherever possible, then exit.")
	offsetPtr := flag.Int("offset", 0, "ignore if 0 or if -limit is not supplied, otherwise use as query OFFSET")
	limitPtr := flag.Int("limit", 0, "ignore if 0, otherwise use as query LIMIT")
	workerPtr := flag.Int("workers", 1, "number of workers to use to send email")
	dbs := flag.Bool("debugsend", false, "print email addresses for recipients but don't send")

	flag.Parse()

	App.Bounce = *bPTR
	App.Complaint = *cPTR
	App.OOO = *oPTR
	App.Suppress = *sPTR
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.MsgFile = *mPtr
	App.AttachFile = *aPtr
	App.QueryName = *qPtr
	App.SetupOnly = *soPtr
	App.Subject = *subjPtr
	App.From = *fromPtr
	App.QueryCount = *qcPtr
	App.MojoHost = *hPtr
	App.Fix = *fixPtr
	App.ValidateGroup = *vPtr
	App.Offset = *offsetPtr
	App.Limit = *limitPtr
	App.DebugSend = *dbs
	App.GroupName = *gPtr
	App.WorkerCount = *workerPtr

	if len(*pids) > 0 {
		sa := strings.Split(*pids, ",")
		for i := 0; i < len(sa); i++ {
			x, err := strconv.ParseInt(sa[i], 10, 64)
			if err != nil {
				panic(err)
			}
			App.PIDs = append(App.PIDs, x)
			fmt.Printf("parsed: %d\n", x)
		}
	}
}

func main() {
	var err error
	fmt.Printf("MOJO Mailsend - begin\n")
	err = db.ReadConfig()
	if err != nil {
		util.UlogAndPrint("Error in db.ReadConfig: %s\n", err.Error())
		os.Exit(1)
	}
	readCommandLineArgs()

	// util.Console("P1 App.QueryName = %s\n", App.QueryName)

	//----------------------------------------------
	// Open the logfile and begin logging...
	//----------------------------------------------
	App.LogFile, err = os.OpenFile("mailsend.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		util.UlogAndPrint("main: %s\n", err.Error())
	}
	defer App.LogFile.Close()
	log.SetOutput(App.LogFile)
	util.Ulog("*** Accord MAILSEND ***\n")

	//----------------------------------------------
	// Open the database...
	//----------------------------------------------
	s := extres.GetSQLOpenString(db.MojoDBConfig.MojoDbname, &db.MojoDBConfig)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		util.UlogAndPrint("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.MojoDBConfig.MojoDbname, db.MojoDBConfig.MojoDbuser, err.Error())
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		util.UlogAndPrint("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", db.MojoDBConfig.MojoDbname, db.MojoDBConfig.MojoDbuser, err.Error())
		os.Exit(1)
	}
	db.InitDB(App.db)
	db.BuildPreparedStatements()

	if App.QueryCount {
		fmt.Printf("Search for query: %s\n", App.QueryName)
		x, err := db.GetQueryRowCount(App.QueryName)
		if err != nil {
			util.UlogAndPrint("Error from GetQueryRowCount: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Row count for query %s = %d\n", App.QueryName, x)
		os.Exit(0)
	}

	// util.Console("P2 App.QueryName = %s\n", App.QueryName)

	si := mailsend.Info{
		From:        App.From,
		QName:       App.QueryName,
		GroupName:   App.GroupName,
		Subject:     App.Subject,
		MsgFName:    App.MsgFile,
		AttachFName: App.AttachFile,
		Hostname:    App.MojoHost,
		SMTPHost:    db.MojoDBConfig.SMTPHost,
		SMTPLogin:   db.MojoDBConfig.SMTPLogin,
		SMTPPass:    db.MojoDBConfig.SMTPPass,
		SMTPPort:    db.MojoDBConfig.SMTPPort,
		Offset:      App.Offset,
		Limit:       App.Limit,
		DebugSend:   App.DebugSend,
		WorkerCount: App.WorkerCount,
	}
	// fmt.Printf("SMTP Info: host:port = %s:%d, login = %s, pass = %s\n", si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
	if App.SetupOnly {
		util.UlogAndPrint("Setup\n")
		setupTestGroups()
		util.UlogAndPrint("Setup completed\n")
		return
	}
	// util.Console("P2.1 App.QueryName = %s, si.QName =%s\n", App.QueryName, si.QName)

	if App.Fix {
		fixDoubleDotEmail()
		fixDotAtEmail()
		return
	}
	// util.Console("P2.2 App.QueryName = %s, si.QName =%s\n", App.QueryName, si.QName)

	if App.Bounce {
		util.UlogAndPrint("Bounce Email\n")
		SendBouncedEmailTest()
		util.UlogAndPrint("Bounce Email Complete\n")
		return
	}
	// util.Console("P2.3 App.QueryName = %s, si.QName =%s\n", App.QueryName, si.QName)
	if App.Complaint {
		util.UlogAndPrint("Complain Email\n")
		SendComplaintEmailTest()
		util.UlogAndPrint("Complain Email Complete\n")
		return
	}
	// util.Console("P2.4 App.QueryName = %s, si.QName =%s\n", App.QueryName, si.QName)
	if App.OOO {
		util.UlogAndPrint("Out-of-Office Email\n")
		SendOOOEmailTest()
		util.UlogAndPrint("Out-of-Office Email Complete\n")
		os.Exit(0)
	}
	// util.Console("P2.5 App.QueryName = %s, si.QName =%s\n", App.QueryName, si.QName)
	if App.Suppress {
		util.UlogAndPrint("Suppression List Email\n")
		SendSuppressionListEmailTest()
		util.UlogAndPrint("Suppression List Email Complete\n")
		os.Exit(0)
	}
	// util.Console("P2.6 App.QueryName = %s, si.QName =%s\n", App.QueryName, si.QName)
	if len(App.ValidateGroup) > 0 {
		fmt.Printf("Validating email addresses for group: %s\n", App.ValidateGroup)
		err = mailsend.ValidateGroupEmailAddresses(App.ValidateGroup)
		if err != nil {
			fmt.Printf("mailsend.ValidateGroupEmailAddresses:  err = %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(App.PIDs) > 0 {
		fmt.Printf("Sending to PID list\n")
		if err = mailsend.SendToPIDs(App.PIDs, &si); err != nil {
			fmt.Printf("Error from SendToPIDs = %s\n", err.Error())
		}
		os.Exit(0)
	}

	// util.Console("P2.7 App.QueryName = %s, si.QName =%s\n", App.QueryName, si.QName)
	if si.Offset > 0 && si.Limit == 0 {
		fmt.Printf("You MUST supply a limit value if you specify an offset\n")
		os.Exit(1)
	}

	// util.Console("P3 App.QueryName = %s, si.QName =%s\n", App.QueryName, si.QName)

	err = mailsend.Sendmail(&si)
	if err != nil {
		util.UlogAndPrint("error sending mail: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully sent %d message(s)\n", si.SentCount)
}
