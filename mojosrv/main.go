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
	"phonebook/lib"
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

// HomeHandler serves static http content such as the .css files
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		Chttp.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/home/", http.StatusFound)
	}
}

// Chttp is a server mux for handling unprocessed html page requests.
// For example, a .css file or an image file.
var Chttp = http.NewServeMux()

func initHTTP() {
	Chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/home/", HomeUIHandler)
	http.HandleFunc("/v1/", V1ServiceHandler)
}

func readCommandLineArgs() {
	portPtr := flag.Int("p", 8275, "port on which mojo server listens")
	flag.Parse()
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
	lib.Errcheck(err)
	defer App.LogFile.Close()
	log.SetOutput(App.LogFile)
	util.Ulog("*** Accord MOJO ***\n")

	// Get the database...
	s := extres.GetSQLOpenString(db.MojoDBConfig.MojoDbname, &db.MojoDBConfig)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.MojoDBConfig.MojoDbname, db.MojoDBConfig.MojoDbuser, err)
		os.Exit(1)
	}
	defer App.db.Close()

	// s := "<awsdbusername>:<password>@tcp(<rdsinstancename>:3306)/accord"
	// s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	// App.db, err = sql.Open("mysql", s)
	// if nil != err {
	// 	fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
	// }
	// defer App.db.Close()

	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", db.MojoDBConfig.MojoDbname, db.MojoDBConfig.MojoDbuser, err)
		os.Exit(1)
	}
	db.InitDB(App.db)
	db.BuildPreparedStatements()
	initHTTP()
	util.Ulog("mojosrv initiating HTTP service on port %d\n", App.Port)

	//go http.ListenAndServeTLS(fmt.Sprintf(":%d", App.Port+1), App.CertFile, App.KeyFile, nil)
	err = http.ListenAndServe(fmt.Sprintf(":%d", App.Port), nil)
	if nil != err {
		fmt.Printf("*** Error on http.ListenAndServe: %v\n", err)
		util.Ulog("*** Error on http.ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
