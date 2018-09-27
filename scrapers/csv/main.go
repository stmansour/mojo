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

// MojoDBFields are the field names into which we can import data
var MojoDBFields = []string{
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
}

// FieldMap describes which imported field goes to which Mojo field
type FieldMap struct {
	MojoField string
	CSVField  string
}

// App is the global application context structure
var App struct {
	db          *sql.DB
	LogFile     *os.File
	DBName      string
	DBUser      string
	fname       string // name of the file we're importing
	debug       bool
	GroupName   string
	CreateGroup bool // if true and group does not exist then create it
	GroupDesc   string
	Group       db.EGroup
	skipOutput  bool // show start/stop time and elapsed time
	keepEmail1  bool // if importing and Email1 != Email2
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

	s := extres.GetSQLOpenString(db.MojoDBConfig.MojoDbname, &db.MojoDBConfig)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if err != nil {
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
		os.Exit(1)
	}

	db.InitDB(App.db)
	db.BuildPreparedStatements()

	//-------------------------
	// Load the group...
	//-------------------------
	if len(App.GroupName) == 0 {
		fmt.Printf("You must supply a group name:  -g groupname\n")
		os.Exit(1)
	}
	App.Group, err = db.GetGroupByName(App.GroupName)
	if nil != err {
		if util.IsSQLNoResultsError(err) {
			if App.CreateGroup {
				var grp db.EGroup
				grp.GroupName = App.GroupName
				grp.DtStart = time.Now()
				grp.GroupDescription = App.GroupDesc
				if err = db.InsertGroup(&grp); err != nil {
					fmt.Printf("Error inserting group %s: %s\n", App.GroupName, err.Error())
					os.Exit(1)
				}
				App.Group = grp
			} else {
				fmt.Printf("Group %s does not exist. Use -cg if you want to create it.\n", App.GroupName)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error opening group %s: %s\n", App.GroupName, err.Error())
			os.Exit(1)
		}
	}

	fmt.Printf("App.fname = %q\n", App.fname)
	if len(App.fname) == 0 {
		fmt.Printf("You must enter -f filename.csv\n")
		os.Exit(1)
	}
	fmt.Printf("               skipOutput = %t\n", App.skipOutput)
	fmt.Printf("                GroupName = %s\n", App.GroupName)
	fmt.Printf("Create group if necessary = %t\n", App.CreateGroup)
	MapAndImport(App.fname)
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "mojo", "database name")
	fPtr := flag.String("f", "", "name of csvfile to parse")
	dbgPtr := flag.Bool("D", false, "use this option to turn on debug mode")
	gPtr := flag.String("g", "", "Add people to this group")
	cgPtr := flag.Bool("cg", false, "Create the group in from -g if necessary")
	gdPtr := flag.String("d", "", "Group description for create (optional)")
	soPtr := flag.Bool("o", false, "Setting this option causes the output to be inhibited")
	flag.Parse()
	App.debug = *dbgPtr
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.fname = *fPtr
	App.GroupName = *gPtr
	App.CreateGroup = *cgPtr
	App.GroupDesc = *gdPtr
	App.skipOutput = *soPtr
}

// MapAndImport looks at the first 2 lines of the csv input file to determine how to map the fields:
//
// line 1 of the csv file should describe all the fields in the schema being imported
// line 2 should have the name of the Mojo field to which it maps, or it can be blank if it does not map
// line 3+ should be the data we're importing
//
// example:
//   First Name,Middle Name,Last Name,License Name, License Type,License Number,License Status,Broker License Number,Broker License Name,License Expires,Email Address,Address Line 1, Address Line 2, City, State, Zip, County, Home Area Code, Home Phone Number,Home Address Line 1, Home Address Line 2, Home City, Home State, Home Zip Code, Home County, DOB
//   FirstName,MiddleName,LastName,,,,,,,,Email1,MailAddress,MailAddress2,MailCity,MailState,MailPostalCode,,,,,,,,,,
//   Lois,Ann,Joyce,Joyce Lois Ann,BR,25630,A,25630,"Joyce, Lois Ann",31-Jan-18,mmjoyce@sbcglobal.net,3432 S GARY
//
// After working out the mapping, it will process lines 3+ and add them to the database
//--------------------------------------------------------------------------------------------------------
func MapAndImport(fname string) {
	var err error
	//-------------------------------------------------------------
	//  First thing to do is establish the mapping between the
	//  columns and mojo's people definition
	//-------------------------------------------------------------
	// util.Console("Entering MapAndImport: fname = %s\n", fname)
	t := util.LoadCSV(fname)
	if len(t) < 3 {
		fmt.Printf("FATAL: only read %d lines from %s\n", len(t), fname)
		os.Exit(1)
	}
	fldmap := []int{}                       // maps the column number of the input csv to the field number of mojo's Person
	inputFields := t[0]                     // we actually just ignore these
	mapToFields := t[1]                     // this one holds the info we need
	for i := 0; i < len(inputFields); i++ { // for each column of the inputfile
		k := -1 // assume we don't find a match
		if len(mapToFields[i]) > 0 {
			for j := 0; j < len(MojoDBFields); j++ {
				if mapToFields[i] == MojoDBFields[j] {
					k = j
					break
				}
			}
		}
		if k == -1 && len(mapToFields[i]) > 0 {
			fmt.Printf("Invalid map-to field name: %s\n", mapToFields[i])
			os.Exit(1)
		}
		fldmap = append(fldmap, k)
	}
	util.Console("fldmap[]:\n")
	count := 0
	for i := 0; i < len(fldmap); i++ {
		util.Console("%d. %d\n", i, fldmap[i])
		if fldmap[i] > 0 {
			count++
		}
	}
	if count == 0 {
		err := fmt.Errorf("No mapping information in second line of csv file could be mapped to a recognized column")
		util.LogAndPrintError("MapAndImport", err)
		os.Exit(1)
	}

	//-------------------------------------------------------------
	// Mark Start time of group update
	//-------------------------------------------------------------
	App.Group, err = db.GetGroupByName(App.GroupName)
	if nil != err {
		log.Fatalf("MapAndImport: error getting group: %s\n", App.GroupName)
	}
	App.Group.DtStart = time.Now()
	err = db.UpdateGroup(&App.Group)
	if nil != err {
		log.Fatalf("MapAndImport: error updating group: %s\n", App.GroupName)
	}
	util.Console("UPDATED App.Group.DtStart = %s\n", App.Group.DtStart.Format(util.DATETIMEINPFMT))

	//-------------------------------------------------------------
	// Now that we know the mapping, go through the data and
	// load the people.
	//-------------------------------------------------------------
	matches := 0
	newEntries := 0
	flt := float64(len(t))
	for i := 2; i < len(t); i++ {
		var a db.Person
		for j := 0; j < len(t[i]); j++ {
			if j >= len(fldmap) {
				continue
			}
			if fldmap[j] == -1 {
				continue
			}
			p := &t[i][j]
			switch fldmap[j] {
			case 0:
				a.FirstName = *p
			case 1:
				a.MiddleName = *p
			case 2:
				a.LastName = *p
			case 3:
				a.PreferredName = *p
			case 4:
				a.JobTitle = *p
			case 5:
				a.OfficePhone = *p
			case 6:
				a.OfficeFax = *p
			case 7:
				a.Email1 = util.ScrubEmailAddr(*p)
			case 8:
				a.Email2 = util.ScrubEmailAddr(*p)
			case 9:
				a.MailAddress = *p
			case 10:
				a.MailAddress2 = *p
			case 11:
				a.MailCity = *p
			case 12:
				a.MailState = *p
			case 13:
				a.MailPostalCode = *p
			case 14:
				a.MailCountry = *p
			case 15:
				a.RoomNumber = *p
			case 16:
				a.MailStop = *p
			default:
				fmt.Printf("unexpected fldmap index: %d\n", fldmap[j])
				os.Exit(1)
			}
		}
		if len(a.Email1) == 0 && len(a.Email2) == 0 {
			// util.Console("no email address found for this entry, skipping to next person\n")
		}
		// util.Console("Processing person:  %q %q (%s)\n", a.FirstName, a.LastName, a.Email1)

		//---------------------------------------------
		// Do we already have this person?
		//---------------------------------------------
		GID := App.Group.GID
		var PID int64
		var dup db.Person
		var err error
		createdPerson := false
		dup, err = db.GetPersonByEmail(a.Email1)
		if err != nil {
			if !util.IsSQLNoResultsError(err) {
				log.Fatalf("Error searching for person with email address %s: %s\n", a.Email1, err.Error())
			}
		}

		//----------------------------------------------------------
		// Add this person if he/she was not found in the database
		//----------------------------------------------------------
		if dup.PID == 0 {
			// util.Console("%s NOT FOUND, adding new\n", a.Email1)
			if err = db.InsertPerson(&a); err != nil { // no: insert the person
				log.Fatalf("Error inserting Person: %s\n", err.Error())
			}
			PID = a.PID // now add this person to the group
			newEntries++
			createdPerson = true //mark that it
		} else {
			matches++
			// util.Console("%s FOUND, PID = %d, will just add to group\n", a.Email1, dup.PID)
			PID = dup.PID // yes: just add the person to this group
		}

		//----------------------------------------------------------
		// Does this person already exist in the group
		//----------------------------------------------------------
		addToGroup := true  // assume we need to add this person to the group
		if !createdPerson { // only need to look if we didn't create the person
			_, err := db.GetPGroup(PID, GID)
			if err != nil {
				if !util.IsSQLNoResultsError(err) {
					log.Fatalf("Error getting person group: %s\n", err.Error())
				} else {
					// util.Console("At point P1 - PID=%d IS NOT already a member of group %d.  addToGroup = %t\n", PID, GID, addToGroup)

				}
			} else {
				addToGroup = false // this person is already a member of the group
				// util.Console("At point P1 - PID=%d IS already a member of group %d.  addToGroup = %t\n", PID, GID, addToGroup)
			}
		}

		//----------------------------------------------------------
		// Add the person to this group if necessary
		//----------------------------------------------------------
		if addToGroup {
			var pg = db.PGroup{
				PID: PID,
				GID: GID,
			}
			// util.Console("Adding %d to group %d\n", PID, GID)
			if err = db.InsertPGroup(&pg); err != nil {
				fmt.Printf("Error inserting Person into group: %s\n", err.Error())
				os.Exit(1)
			}
		}

		//----------------------------------------------------------
		// update on-screen progress
		//----------------------------------------------------------
		if i%100 == 0 {
			fmt.Printf("\r%8d  -->  %3.1f%%", i, float64(100*i)/flt)
		}
	}

	//-------------------------------------------------------------
	// Mark Stop time of group update
	//-------------------------------------------------------------
	App.Group, err = db.GetGroupByName(App.GroupName)
	if nil != err {
		log.Fatalf("MapAndImport: error getting group: %s\n", App.GroupName)
	}
	App.Group.DtStop = time.Now()
	err = db.UpdateGroup(&App.Group)
	if nil != err {
		log.Fatalf("MapAndImport: error updating group: %s\n", App.GroupName)
	}
	util.Console("UPDATED App.Group.DtStop = %s\n", App.Group.DtStop.Format(util.DATETIMEINPFMT))

	//-------------------------------------------------------------
	// Print out the stats...
	//-------------------------------------------------------------
	fmt.Printf("Import Complete\n")
	if !App.skipOutput {
		fmt.Printf("Start time:    %s\n", App.Group.DtStart.In(db.DB.Zone).Format(util.DATETIMEINPFMT))
		fmt.Printf("Stop time:     %s\n", App.Group.DtStop.In(db.DB.Zone).Format(util.DATETIMEINPFMT))
		fmt.Printf("Elapsed time:  %s\n", App.Group.DtStop.Sub(App.Group.DtStart))
		fmt.Printf("Matched:       %6d\n", matches)
		fmt.Printf("New Entries:   %6d\n", newEntries)
		fmt.Printf("Total Updates: %6d\n", matches+newEntries)
	}
}
