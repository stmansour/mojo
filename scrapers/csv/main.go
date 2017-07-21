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
	flag.Parse()
	App.debug = *dbgPtr
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.fname = *fPtr
	App.GroupName = *gPtr
	App.CreateGroup = *cgPtr
	App.GroupDesc = *gdPtr
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
	//-------------------------------------------------------------
	//  First thing to do is establish the mapping between the
	//  columns and mojo's people definition
	//-------------------------------------------------------------
	t := util.LoadCSV(fname)
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
	// fmt.Printf("fldmap[]:\n")
	// for i := 0; i < len(fldmap); i++ {
	// 	fmt.Printf("%d. %d\n", i, fldmap[i])
	// }

	//-------------------------------------------------------------
	// Now that we know the mapping, go through the data anc
	// load the people.
	//-------------------------------------------------------------
	flt := float64(len(t))
	for i := 2; i < len(t); i++ {
		var a db.Person
		for j := 0; j < len(t[i]); j++ {
			if fldmap[j] == -1 {
				continue
			}
			p := &t[i][j]
			switch fldmap[j] {
			case 0: // FirstName
				a.FirstName = *p
			case 1: // MiddleName
				a.MiddleName = *p
			case 2: // LastName
				a.LastName = *p
			case 3: // PreferredName
				a.PreferredName = *p
			case 4: // JobTitle
				a.JobTitle = *p
			case 5: // OfficePhone
				a.OfficePhone = *p
			case 6: // OfficeFax
				a.OfficeFax = *p
			case 7: // Email1
				a.Email1 = *p
			case 8: // MailAddress
				a.MailAddress = *p
			case 9: // MailAddress2
				a.MailAddress2 = *p
			case 10: // MailCity
				a.MailCity = *p
			case 11: // MailState
				a.MailState = *p
			case 12: // MailPostalCode
				a.MailPostalCode = *p
			case 13: // MailCountry
				a.MailCountry = *p
			case 14: // RoomNumber
				a.RoomNumber = *p
			case 15: // MailStop
				a.MailStop = *p
			default:
				fmt.Printf("unexpected fldmap index: %d\n", fldmap[j])
				os.Exit(1)
			}
		}

		// fmt.Printf("Adding person:  %s %s (%s)\n", a.FirstName, a.LastName, a.Email1)

		// Do we already have this person?
		GID := App.Group.GID
		var PID int64
		var dup db.Person
		var err error
		dup, err = db.GetPersonByEmail(a.Email1)
		if err != nil {
			if !util.IsSQLNoResultsError(err) {
				fmt.Printf("Error searching for person with email address %s: %s\n", a.Email1, err.Error())
				os.Exit(1)
			}
		}
		if dup.PID == 0 { // did we find the person...
			if err = db.InsertPerson(&a); err != nil { // no: insert the person
				fmt.Printf("Error inserting Person: %s\n", err.Error())
				os.Exit(1)
			}
			PID = a.PID // now add this person to the group
		} else {
			PID = dup.PID // yes: just add the person to this group
		}
		var pg = db.PGroup{
			PID: PID,
			GID: GID,
		}
		if err = db.InsertPGroup(&pg); err != nil {
			fmt.Printf("Error inserting Person into group: %s\n", err.Error())
			os.Exit(1)
		}

		if i%100 == 0 {
			fmt.Printf("\r%8d  -->  %3.1f%%", i, float64(100*i)/flt)
		}

	}
}
