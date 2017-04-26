package main

import (
	"extres"
	"flag"
	"fmt"
	"html/template"
	"mojo/db"
	"mojo/mailsend"
	"mojo/util"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// App is the application structure available to the whole app.
var App struct {
	From     string
	MsgFile  string
	MojoHost string // domain and port for mojosrv:   http://example.domain.com:8275/
	Subject  string
	XR       extres.ExternalResources // dbs, smtp...
}

func readCommandLineArgs() {
	mPtr := flag.String("b", "testmsg.html", "filename containing the html message to send")
	hPtr := flag.String("h", "http://localhost:8275/", "name of host and port for mojosrv")
	subjPtr := flag.String("subject", "Test Message", "Email subject line.")
	fromPtr := flag.String("from", "sman@accordinterests.com", "Message sender.")

	flag.Parse()
	App.MsgFile = *mPtr
	App.Subject = *subjPtr
	App.From = *fromPtr
	App.MojoHost = *hPtr
}

func main() {
	readCommandLineArgs()
	var err error
	err = db.ReadConfig()
	if err != nil {
		util.UlogAndPrint("Error in db.ReadConfig: %s\n", err.Error())
		os.Exit(1)
	}

	t, err := template.New(App.MsgFile).ParseFiles(App.MsgFile)
	if nil != err {
		fmt.Printf("error loading template: %v\n", err)
		os.Exit(1)
	}

	var p = db.Person{FirstName: "Steven", MiddleName: "F", LastName: "Mansour", JobTitle: "CTO, Accord Interests", OfficePhone: "323-512-0111 X305", Email1: "sman@accordinterests.com", MailAddress: "11719 Bee Cave Road", MailAddress2: "Suite 301", MailCity: "Austin", MailState: "TX", MailPostalCode: "78738", MailCountry: "USA", Status: 0}
	th, err := mailsend.GeneratePageHTML(App.MsgFile, App.MojoHost, &p, t)
	fmt.Printf("Body string:\n\n%s\n", string(th))
}
