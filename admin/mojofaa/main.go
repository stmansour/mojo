// faa  a program to scrape the FAA directory site.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"mojo/db"
	"mojo/scrapers/faa"
	"mojo/util"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global data structure for this app
var App struct {
	db        *sql.DB
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
			err = db.InsertGroup(&g)
			if err != nil {
				fmt.Printf("Error inserting group: %s\n", err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error reading group \"FAA\": %s\n", err.Error())
			os.Exit(1)
		}
	}
	App.GID = g.GID

}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "mojo", "database name")
	fPtr := flag.String("f", "step3.csv", "name of csvfile to parse")
	sPtr := flag.String("s", "", "skip names until you find this name, then engage")
	dbgPtr := flag.Bool("D", false, "use this option to turn on debug mode")
	qPtr := flag.Bool("q", false, "quick mode - only loop once - enables fast start to finish testing")
	wpPtr := flag.Int("w", 25, "Number of workers in the worker pool")
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
	db.BuildPreparedStatements()
	faa.InitFAAScraper(App.GID, App.quick, App.workers)
	faa.ScrapeFAA()
}
