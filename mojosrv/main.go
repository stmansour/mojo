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
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global data structure for this app
var App struct {
	db        *sql.DB
	DBName    string
	DBUser    string
	Port      int      // port on which mojo listens
	LogFile   *os.File // where to log messages
	fname     string
	startName string
}

// Chttp is a server mux for handling unprocessed html page requests.
// For example, a .css file or an image file.
var Chttp = http.NewServeMux()

func unauthorizeFileReqHandler(w http.ResponseWriter, r *http.Request) {
	util.Ulog("UNAUTHORIZED FILE ACCESS ATTEMPT: url = %s\n", r.URL.Path)
	fmt.Fprintf(w, "404 page not found")
}

// HomeHandler serves static http content such as the .css files
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		Chttp.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/home/", http.StatusFound)
	}
}

func initHTTP() {
	http.HandleFunc("/config.json", unauthorizeFileReqHandler)
	http.HandleFunc("/confprod.json", unauthorizeFileReqHandler)
	http.HandleFunc("/confdev.json", unauthorizeFileReqHandler)
	Chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/home/", HomeUIHandler)
	http.HandleFunc("/v1/", V1ServiceHandler)
}

func readCommandLineArgs() {
	portPtr := flag.Int("p", 8275, "port on which mojo server listens")
	vptr := flag.Bool("v", false, "Show version, then exit")
	flag.Parse()
	if *vptr {
		fmt.Printf("Version:   %s\n", getVersionNo())
		os.Exit(0)
	}
	App.Port = *portPtr
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
	App.LogFile, err = os.OpenFile("mojo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	util.ErrCheck(err)
	defer App.LogFile.Close()
	log.SetOutput(App.LogFile)
	util.Ulog("*** Accord MOJO ***\n")

	// Get the database...
	// s := "<awsdbusername>:<password>@tcp(<rdsinstancename>:3306)/accord"
	s := extres.GetSQLOpenString(db.MojoDBConfig.MojoDbname, &db.MojoDBConfig)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.MojoDBConfig.MojoDbname, db.MojoDBConfig.MojoDbuser, err)
		os.Exit(1)
	}
	defer App.db.Close()

	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", db.MojoDBConfig.MojoDbname, db.MojoDBConfig.MojoDbuser, err)
		os.Exit(1)
	}
	db.InitDB(App.db)
	db.BuildPreparedStatements()
	initHTTP()
	util.Ulog("mojosrv initiating HTTP service on port %d\n", App.Port)
	util.Ulog("Using database: %s , host = %s, port = %d\n", db.MojoDBConfig.MojoDbname, db.MojoDBConfig.MojoDbhost, db.MojoDBConfig.MojoDbport)
	fmt.Printf("Using database: %s , host = %s, port = %d\n", db.MojoDBConfig.MojoDbname, db.MojoDBConfig.MojoDbhost, db.MojoDBConfig.MojoDbport)

	//go http.ListenAndServeTLS(fmt.Sprintf(":%d", App.Port+1), App.CertFile, App.KeyFile, nil)
	err = http.ListenAndServe(fmt.Sprintf(":%d", App.Port), nil)
	if nil != err {
		fmt.Printf("*** Error on http.ListenAndServe: %v\n", err)
		util.Ulog("*** Error on http.ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
