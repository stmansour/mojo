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
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// App is the application structure available to the whole app.
var App struct {
	GroupName string
	db        *sql.DB
	DBName    string
	DBUser    string
	LogFile   *os.File // where to log messages
	// XR        extres.ExternalResources // dbs, smtp...
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "mojo", "database name")
	gPtr := flag.String("group", "", "group name to export")

	flag.Parse()

	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.GroupName = *gPtr

}

func main() {
	var err error
	err = db.ReadConfig()
	if err != nil {
		util.UlogAndPrint("Error in db.ReadConfig: %s\n", err.Error())
		os.Exit(1)
	}
	readCommandLineArgs()

	//----------------------------------------------
	// Open the logfile and begin logging...
	//----------------------------------------------
	App.LogFile, err = os.OpenFile("mojoexport.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		util.UlogAndPrint("main: %s\n", err.Error())
	}
	defer App.LogFile.Close()
	log.SetOutput(App.LogFile)
	util.Ulog("*** Accord MOJO-EXPORT ***\n")

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

	doExport()
}

func doExport() {
	var err error
	var GID = int64(0)
	var joins string
	var egrp db.EGroup

	if len(App.GroupName) > 0 {
		// util.Console("GROUP NAME = %s\n", App.GroupName)

		//---------------------------------------------------------------------
		// if the group name is valid, we want a query like this one:
		// SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=GID
		//---------------------------------------------------------------------
		egrp, err = db.GetGroupByName(App.GroupName)
		if err != nil {
			util.Console("Error from db.GetGroupByName: %s\n", err.Error())
			return
		}
		GID = egrp.GID
		util.Console("GID = %d\n", GID)
	}

	s1 := db.DB.DBFields["People"] // comma separated list
	sa := strings.Split(s1, ",")
	for i := 0; i < len(sa); i++ {
		sa[i] = "People." + sa[i] // remove any ambiguity after the join
	}
	flds := strings.Join(sa, ",")

	q := fmt.Sprintf("SELECT %s FROM People ", flds) // the fields we want
	if GID > 0 {
		joins = fmt.Sprintf("INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d ", GID)
		q += joins
	}

	q += " ORDER BY LastName ASC, Firstname ASC, MiddleName ASC"

	// util.Console("\nQuery = %s\n\n", q)

	rows, err := db.DB.Db.Query(q)
	if err != nil {
		util.Console("Error from DB Query: %s\n", err.Error())
		return
	}
	defer rows.Close()

	exportHeader()

	count := 0
	for rows.Next() {
		p, err := db.ReadPersonFromRows(rows)
		if err != nil {
			util.Console("%s.  Error reading Person: %s\n", "doExport", err.Error())
		}
		exportHandle(p)
		count++ // update the count only after adding the record
	}
	// util.Console("Exported Records:  %d\n", count)
	return
}

func exportHandle(p db.Person) {
	s := fmt.Sprintf("%d,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%d,%s\n",
		p.PID,
		p.FirstName,
		p.MiddleName,
		p.LastName,
		p.PreferredName,
		p.JobTitle,
		p.OfficePhone,
		p.OfficeFax,
		p.Email1,
		p.Email2,
		p.MailAddress,
		p.MailAddress2,
		p.MailCity,
		p.MailState,
		p.MailPostalCode,
		p.MailCountry,
		p.RoomNumber,
		p.MailStop,
		p.Status,
		p.OptOutDate.Format(util.DATEINPFMT),
	)
	fmt.Print(s)
}

func exportHeader() {
	var cols = []string{
		"PID",
		"FirstName",
		"MiddleName",
		"LastName",
		"PreferredName",
		"JobTitle",
		"OfficePhone",
		"OfficeFax",
		"Email1",
		"Email2",
		"MailAddress",
		"MailAddress2",
		"MailCity",
		"MailState",
		"MailPostalCode",
		"MailCountry",
		"RoomNumber",
		"MailStop",
		"Status",
		"OptOutDate",
	}
	for i := 0; i < len(cols); i++ {
		fmt.Printf("%q", cols[i])
		if i != len(cols)-1 {
			fmt.Printf(",")
		}
	}
	fmt.Println()

}
