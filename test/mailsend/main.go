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
	"time"

	"gopkg.in/gomail.v2"

	_ "github.com/go-sql-driver/mysql"
)

// App is the application structure available to the whole app.
var App struct {
	MsgFile    string
	AttachFile string
	QueryName  string
	MojoHost   string // domain and port for mojosrv:   http://example.domain.com:8275/
	db         *sql.DB
	DBName     string
	DBUser     string
	SetupOnly  bool
	Subject    string                   // subject line of the email message
	From       string                   // email from address
	QueryCount bool                     // if true, just print the solution set count for the query and exit
	Bounce     bool                     // if true, just print the solution set count for the query and exit
	Complaint  bool                     // if true, just print the solution set count for the query and exit
	OOO        bool                     // if true, just print the solution set count for the query and exit
	Suppress   bool                     // if true, just print the solution set count for the query and exit
	LogFile    *os.File                 // where to log messages
	XR         extres.ExternalResources // dbs, smtp...
}

// SendBouncedEmailTest sends an email message that bounces.  For testing.
// The recipient's ISP rejects your email with an SMTP 550 5.1.1 response
// code ("Unknown User"). Amazon SES generates a bounce notification and
// sends it to you via email or by using an Amazon SNS notification,
// depending on how you set up your system. This mailbox simulator email
// address will not be placed on the Amazon SES suppression list as one
// normally would when an email hard bounces. The bounce response that
// you receive from the mailbox simulator is compliant with RFC 3464.
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
func SendOOOEmailTest() error {
	return SendEmailTest("ooto@simulator.amazonses.com")
}

// SendSuppressionListEmailTest sends a test email. Amazon SES treats your
// email as a hard bounce because the address you are sending to is on the
// Amazon SES suppression list.
func SendSuppressionListEmailTest() error {
	return SendEmailTest("suppressionlist@simulator.amazonses.com")
}

// SendEmailTest is a routine to send an email message to the supplied address.
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

func addPerson(pnew *db.Person, GID int64) error {
	var pid int64
	p1, err := db.GetPersonByName(pnew.FirstName, pnew.MiddleName, pnew.LastName)
	if err != nil {
		util.Ulog("db.GetPersonByName returned: %s\n", err.Error())
		return err
	}
	if len(p1) == 0 {
		err := db.InsertPerson(pnew)
		if err != nil {
			util.Ulog("db.InsertPerson returned: %s\n", err.Error())
			return err
		}
		pid = pnew.PID
	} else {
		pid = p1[0].PID
	}
	if GID == int64(0) {
		return nil
	}
	return AddPersonToGroup(pid, GID)
}

func resetUserStatus(ppa *[]db.Person) {
	pa := *ppa
	for i := 0; i < len(pa); i++ {
		p, err := db.GetPersonByEmail(pa[i].Email1)
		if err != nil {
			util.UlogAndPrint("Error from db.GetPersonByEmail( %s ):  %s \n", pa[i].Email1, err.Error())
			os.Exit(1)
		}
		p.Status = db.NORMAL
		if err = db.UpdatePerson(&p); err != nil {
			util.UlogAndPrint("Error from db.GetPersonByEmail( %s ):  %s \n", pa[i].Email1, err.Error())
			os.Exit(1)
		}
	}
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
		resetUserStatus(ppa)
		return
	}
	if err != nil {
		if !util.IsSQLNoResultsError(err) {
			util.UlogAndPrint("Error reading group \"MojoTest\": %s\n", err.Error())
			os.Exit(1)
		}
	}

	// Create the group...
	g.GroupName = name
	g.GroupDescription = descr
	g.DtStart = time.Now()
	err = db.InsertGroup(&g)
	if err != nil {
		util.UlogAndPrint("Error inserting group: %s\n", err.Error())
		os.Exit(1)
	}

	// Add the list of people to it...
	pa := *ppa
	for i := 0; i < len(pa); i++ {
		gid := g.GID
		addPerson(&pa[i], gid)
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
	}
}

func createFAAQueries() {
	g, err := db.GetGroupByName("FAA")
	if err != nil {
		util.UlogAndPrint("Error getting group FAA: %s\n", err.Error())
		os.Exit(1)
	}
	q := fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0 LIMIT 50 OFFSET 0", g.GID) // 50
	createQuery("FAA-1-First50", "The first 50 people in the FAA", q)

	q = fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0 LIMIT 250 OFFSET 50", g.GID) // 300
	createQuery("FAA-2-Next250", "After FAA-1-First50, the next 250 people in the FAA", q)

	q = fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0 LIMIT 700 OFFSET 300", g.GID) // 1,000
	createQuery("FAA-3-Next700", "After FAA-2-Next250, the next 700 people in the FAA", q)

	q = fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0 LIMIT 5000 OFFSET 1000", g.GID) // 6,000
	createQuery("FAA-4-Next5000", "After FAA-3-Next700, the next 700 people in the FAA", q)

	q = fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0 LIMIT 20000 OFFSET 6000", g.GID) // 26,000
	createQuery("FAA-5-Next20000", "After FAA-4-Next5000, the next 20000 people in the FAA", q)

	q = fmt.Sprintf("SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0 LIMIT 50000 OFFSET 26000", g.GID) // up to 56,0000 people
	createQuery("FAA-6-TheRest", "After FAA-5-Next20000, the remaining people in the FAA", q)
}

func setupTestGroups() {
	var pa = []db.Person{
		{FirstName: "Steven", MiddleName: "F", LastName: "Mansour", JobTitle: "CTO, Accord Interests", OfficePhone: "323-512-0111 X305", Email1: "sman@accordinterests.com", MailAddress: "11719 Bee Cave Road", MailAddress2: "Suite 301", MailCity: "Austin", MailState: "TX", MailPostalCode: "78738", MailCountry: "USA", Status: 0},
		{FirstName: "Steve", MiddleName: "", LastName: "Mansour", JobTitle: "Recording Musician, Engineer, Producer", OfficePhone: "323-512-0111 X305", Email1: "sman@stevemansour.com", MailAddress: "2215 Wellington Drive", MailAddress2: "", MailCity: "Milpitas", MailState: "CA", MailPostalCode: "95035", MailCountry: "USA", Status: 0},
	}
	createGroup("MojoTest", "Steve-only test group", &pa)

	var pa1 = []db.Person{
		pa[0],
		pa[1],
		{FirstName: "Bouncie", MiddleName: "", LastName: "McBounce", JobTitle: "Vagabond", OfficePhone: "123-456-7890", Email1: "bounce@simulator.amazonses.com", MailAddress: "123 Elm St", MailAddress2: "", MailCity: "Anytown", MailState: "CA", MailPostalCode: "90210", MailCountry: "USA", Status: 0},
		{FirstName: "Wendy", MiddleName: "", LastName: "Whiner", JobTitle: "Complainer", OfficePhone: "123-321-7890", Email1: "complaint@simulator.amazonses.com", MailAddress: "321 Elm St", MailAddress2: "", MailCity: "Anytown", MailState: "CA", MailPostalCode: "90210", MailCountry: "USA", Status: 0},
		{FirstName: "Stealthy", MiddleName: "", LastName: "McStealth", JobTitle: "Bad Guy", OfficePhone: "816-321-0123", Email1: "suppressionlist@simulator.amazonses.com", MailAddress: "700 Elm St", MailAddress2: "", MailCity: "Anytown", MailState: "CA", MailPostalCode: "90210", MailCountry: "USA", Status: 0},
	}
	createGroup("AmazonTest", "Steve + Amazon test accounts", &pa1)

	var pa2 = []db.Person{
		pa1[0],
		pa1[1],
		pa1[2],
		pa1[3],
		pa1[4],
		{FirstName: "Joe", MiddleName: "G", LastName: "Mansour", JobTitle: "Principal, Accord Interests", OfficePhone: "323-512-0111 X303", Email1: "jgm@accordinterests.com", MailAddress: "11719 Bee Cave Road", MailAddress2: "Suite 301", MailCity: "Austin", MailState: "TX", MailPostalCode: "78738", MailCountry: "USA", Status: 0},
		{FirstName: "Melissa", MiddleName: "", LastName: "Wheeler", JobTitle: "General Manager, Isola Bella", OfficePhone: "405.721.2194 x205", Email1: "mwheeler@myisolabella.com", MailAddress: "8309 NW 140th St", MailAddress2: "", MailCity: "Oklahoma City", MailState: "OK", MailPostalCode: "73142", MailCountry: "USA", Status: 0},
	}
	createGroup("AccordTest", "Steve + Amazon test + Accord accounts", &pa2)
	createFAAQueries()
}
func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "mojo", "database name")
	mPtr := flag.String("b", "testmsg.html", "filename containing the html message to send")
	aPtr := flag.String("a", "", "filename of attachment")
	qPtr := flag.String("q", "MojoTest", "name of the query to send messages to")
	hPtr := flag.String("h", db.MojoDBConfig.MojoWebAddr, "name of host and port for mojosrv")
	qcPtr := flag.Bool("count", false, "returns the count of target addresses in the query, then exits.")
	soPtr := flag.Bool("setup", false, "just run the setup, do not send email")
	bPTR := flag.Bool("bounce", false, "just send a message to bounce@simulator.amazonses.com")
	cPTR := flag.Bool("complaint", false, "just send a message to complaint@simulator.amazonses.com")
	oPTR := flag.Bool("ooo", false, "just send a message to ooo@simulator.amazonses.com")
	sPTR := flag.Bool("sl", false, "just send a message to suppressionlist@simulator.amazonses.com")
	subjPtr := flag.String("subject", "Test Message", "Email subject line.")
	fromPtr := flag.String("from", "sman@accordinterests.com", "Message sender.")

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
	//----------------------------------------------
	// Open the logfile and begin logging...
	//----------------------------------------------
	App.LogFile, err = os.OpenFile("mailsend.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		util.UlogAndPrint("main", err.Error())
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

	si := mailsend.Info{
		From:        App.From,
		QName:       App.QueryName,
		Subject:     App.Subject,
		MsgFName:    App.MsgFile,
		AttachFName: App.AttachFile,
		Hostname:    App.MojoHost,
		SMTPHost:    db.MojoDBConfig.SMTPHost,
		SMTPLogin:   db.MojoDBConfig.SMTPLogin,
		SMTPPass:    db.MojoDBConfig.SMTPPass,
		SMTPPort:    db.MojoDBConfig.SMTPPort,
	}
	// fmt.Printf("SMTP Info: host:port = %s:%d, login = %s, pass = %s\n", si.SMTPHost, si.SMTPPort, si.SMTPLogin, si.SMTPPass)
	if App.SetupOnly {
		util.UlogAndPrint("Setup\n")
		setupTestGroups()
		util.UlogAndPrint("Setup completed\n")
		return
	}
	if App.Bounce {
		util.UlogAndPrint("Bounce Email\n")
		SendBouncedEmailTest()
		util.UlogAndPrint("Bounce Email Complete\n")
		return
	}
	if App.Complaint {
		util.UlogAndPrint("Complain Email\n")
		SendComplaintEmailTest()
		util.UlogAndPrint("Complain Email Complete\n")
		return
	}
	if App.OOO {
		util.UlogAndPrint("Out-of-Office Email\n")
		SendOOOEmailTest()
		util.UlogAndPrint("Out-of-Office Email Complete\n")
		os.Exit(0)
	}
	if App.Suppress {
		util.UlogAndPrint("Suppression List Email\n")
		SendSuppressionListEmailTest()
		util.UlogAndPrint("Suppression List Email Complete\n")
		os.Exit(0)
	}

	err = mailsend.Sendmail(&si)
	if err != nil {
		util.UlogAndPrint("error sending mail: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully sent %d message(s)\n", si.SentCount)
}
