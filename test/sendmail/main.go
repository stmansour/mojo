package main

import (
	"database/sql"
	"flag"
	"fmt"
	"mojo/db"
	"mojo/sendmail"
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
	db         *sql.DB
	DBName     string
	DBUser     string
	SetupOnly  bool
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
	d := gomail.NewDialer("email-smtp.us-east-1.amazonaws.com", 587, "AKIAJ3PENIYLS5U5ATJA", "AqIWufI4PwuxA61NihNQ4Yt+23n6w0CuQLuiUAdHP2E7")
	err := d.DialAndSend(m)
	if err != nil {
		util.Ulog("Error on DialAndSend = %s\n", err.Error())
		return err
	}
	fmt.Printf("Bount message successfully sent to %s\n", addr)
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

func addPerson(p *db.Person, GID int64) error {
	err := db.InsertPerson(p)
	if err != nil {
		util.Ulog("db.InsertPerson returned: %s\n", err.Error())
		return err
	}
	return AddPersonToGroup(p.PID, GID)
}

func setupTestGroup() {
	var g db.EGroup
	var err error
	g, err = db.GetGroupByName("MojoTest")
	if err != nil {
		if util.IsSQLNoResultsError(err) {
			g.GroupName = "MojoTest"
			g.GroupDescription = "Steve's test group"
			g.DtStart = time.Now()
			err = db.InsertGroup(&g)
			if err != nil {
				fmt.Printf("Error inserting group: %s\n", err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error reading group \"MojoTest\": %s\n", err.Error())
			os.Exit(1)
		}
		var p db.Person
		p.FirstName = "Steve"
		p.MiddleName = "F"
		p.LastName = "Mansour"
		p.PreferredName = "Steve"
		p.JobTitle = "CTO, Accord Interests"
		p.OfficePhone = "323-512-0111 X305"
		p.OfficeFax = ""
		p.Email1 = "sman@accordinterests.com"
		p.MailAddress = "11719 Bee Cave Road"
		p.MailAddress2 = "Suite 301"
		p.MailCity = "Austin"
		p.MailState = "TX"
		p.MailPostalCode = "78738"
		p.MailCountry = "USA"
		p.RoomNumber = ""
		p.MailStop = ""
		p.Status = 1 // mark as opt out
		addPerson(&p, g.GID)

		p.JobTitle = "Recording Musician, Engineer, Producer"
		p.Email1 = "sman@stevemansour.com"
		p.MailAddress = "2215 Wellington Drive"
		p.MailAddress2 = ""
		p.MailCity = "Milpitas"
		p.MailState = "CA"
		p.MailPostalCode = "95035"
		p.Status = 0 // normal
		addPerson(&p, g.GID)
	} else {
		// if it's already in the database, we update the record to force the
		// last modified date to reflect the fact that we're scraping now
		fmt.Printf("MojoTest exists, updating timestamp\n")
		g.DtStart = time.Now()
		err = db.UpdateGroup(&g)
		if err != nil {
			fmt.Printf("Error updating group: %s\n", err.Error())
			os.Exit(1)
		}
	}

	var q db.Query
	q, err = db.GetQueryByName("MojoTest")
	if err != nil {
		if util.IsSQLNoResultsError(err) {
			q.QueryName = "MojoTest"
			q.QueryDescr = "Steve's test query"

			// TBD: until we work out the generic sql query builder,
			// I will just store the actual query for now.  This will be
			// replaced when the query builder is completed
			q.QueryJSON = "SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=2 WHERE People.Status=0"
			err = db.InsertQuery(&q)
			if err != nil {
				fmt.Printf("Error inserting query: %s\n", err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error reading query \"MojoTest\": %s\n", err.Error())
			os.Exit(1)
		}
	}
	App.QueryName = q.QueryName
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "mojo", "database name")
	mPtr := flag.String("b", "msg.html", "filename containing the html message to send")
	aPtr := flag.String("a", "", "filename of attachment")
	qPtr := flag.String("q", "MojoTest", "name of the query to send messages to")
	soPtr := flag.Bool("n", false, "just run the setup, do not send email")
	bPTR := flag.Bool("bounce", false, "just send a message to bounce@simulator.amazonses.com")
	cPTR := flag.Bool("complaint", false, "just send a message to complaint@simulator.amazonses.com")
	oPTR := flag.Bool("ooo", false, "just send a message to ooo@simulator.amazonses.com")
	sPTR := flag.Bool("sl", false, "just send a message to suppressionlist@simulator.amazonses.com")

	flag.Parse()
	if *bPTR {
		SendBouncedEmailTest()
		os.Exit(0)
	}
	if *cPTR {
		SendComplaintEmailTest()
		os.Exit(0)
	}
	if *oPTR {
		SendOOOEmailTest()
		os.Exit(0)
	}
	if *sPTR {
		SendSuppressionListEmailTest()
		os.Exit(0)
	}
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.MsgFile = *mPtr
	App.AttachFile = *aPtr
	App.QueryName = *qPtr
	App.SetupOnly = *soPtr
}

func main() {
	readCommandLineArgs()

	var err error
	// s := "<awsdbusername>:<password>@tcp(<rdsinstancename>:3306)/accord"
	s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
		os.Exit(1)
	}
	db.InitDB(App.db)
	db.BuildPreparedStatements()

	si := sendmail.Info{
		From:        "sman@accordinterests.com",
		QName:       App.QueryName,
		Subject:     "Perks Email",
		MsgFName:    App.MsgFile,
		AttachFName: App.AttachFile,
	}
	setupTestGroup()
	if App.SetupOnly {
		fmt.Printf("Setup completed\n")
	} else {
		err = sendmail.Sendmail(&si)
		if err != nil {
			fmt.Printf("error sending mail: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Successfully sent %d message(s)\n", si.SentCount)
	}
}
