// faa  a program to scrape the FAA directory site.
package main

import (
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"log"
	"mojo/db"
	"mojo/util"
	"os"
	"phonebook/lib"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global data structure for this app
var App struct {
	db        *sql.DB
	LogFile   *os.File
	DBName    string
	DBUser    string
	fname     string
	startName string
	debug     bool
	workers   int // number of workers in the goroutine worker pool
	quick     bool
	GID       int64 // GID of the FAA group
}

func initFAAMojo() {
	var g db.EGroup
	var err error
	g, err = db.GetGroupByName("FAA")
	if err != nil {
		if util.IsSQLNoResultsError(err) {
			g.GroupName = "FAA"
			g.GroupDescription = "Employees of the Federal Aviation Administration"
			g.DtStart = time.Now()
			err = db.InsertGroup(&g)
			if err != nil {
				fmt.Printf("Error inserting group: %s\n", err.Error())
				os.Exit(1)
			}
			App.GID = g.GID
			var p db.Person
			p.FirstName = "Steve"
			p.MiddleName = "F"
			p.LastName = "Mansour"
			p.PreferredName = "Steve"
			p.JobTitle = "CTO, Accord Interests"
			p.OfficePhone = ""
			p.OfficeFax = ""
			p.Email1 = "sman@accordinterests.com"
			p.MailAddress = "2215 Wellington Drive"
			p.MailAddress2 = ""
			p.MailCity = "Milpitas"
			p.MailState = "CA"
			p.MailPostalCode = "95035"
			p.MailCountry = "USA"
			p.RoomNumber = ""
			p.MailStop = ""
		} else {
			fmt.Printf("Error reading group \"FAA\": %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		// if it's already in the database, we update the record to force the
		// last modified date to reflect the fact that we're scraping now
		fmt.Printf("FAA exists, updating timestamp\n")
		g.DtStart = time.Now()
		App.GID = g.GID
		err = db.UpdateGroup(&g)
		if err != nil {
			fmt.Printf("Error updating group: %s\n", err.Error())
			os.Exit(1)
		}
	}

}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "mojo", "database name")
	fPtr := flag.String("f", "step3.csv", "name of csvfile to parse")
	sPtr := flag.String("s", "", "skip names until you find this name, then engage")
	dbgPtr := flag.Bool("D", false, "use this option to turn on debug mode")
	qPtr := flag.Bool("q", false, "quick mode - only loop once - enables fast start to finish testing")
	wpPtr := flag.Int("w", 35, "Number of workers in the worker pool")
	flag.Parse()
	App.debug = *dbgPtr
	App.workers = *wpPtr
	App.quick = *qPtr
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.fname = *fPtr
	App.startName = *sPtr
}

func main() {
	var err error
	readCommandLineArgs()
	err = db.ReadConfig()
	if err != nil {
		fmt.Printf("Error in db.ReadConfig: %s\n", err.Error())
		os.Exit(1)
	}

	//==============================================
	// Open the logfile and begin logging...
	//==============================================
	App.LogFile, err = os.OpenFile("scrapefaa.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	lib.Errcheck(err)
	defer App.LogFile.Close()
	log.SetOutput(App.LogFile)
	util.Ulog("*** Accord MOJO FAA Scraper ***\n")

	// s := "<awsdbusername>:<password>@tcp(<rdsinstancename>:3306)/accord"
	// s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	s := extres.GetSQLOpenString(db.MojoDBConfig.MojoDbname, &db.MojoDBConfig)
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
	initFAAMojo()
	if App.GID == int64(0) {
		fmt.Printf("App.GID == 0.  Something is wrong.\n")
		os.Exit(1)
	}
	InitFAAScraper(App.GID, App.quick, App.workers)
	ScrapeFAA()
}
