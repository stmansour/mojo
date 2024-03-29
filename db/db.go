package db

import (
	"context"
	"database/sql"
	"fmt"
	"mojo/util"
	"strings"
	"time"
)

// Person is a structure of all attributes of the FAA employees we're capturing
// Person is the structure that defines all the attributes of a person
type Person struct {
	PID            int64
	FirstName      string
	MiddleName     string
	LastName       string
	PreferredName  string
	JobTitle       string
	OfficePhone    string
	OfficeFax      string
	Email1         string
	Email2         string
	MailAddress    string
	MailAddress2   string
	MailCity       string
	MailState      string
	MailPostalCode string
	MailCountry    string
	RoomNumber     string
	MailStop       string
	Status         int64
	OptOutDate     time.Time
	LastModTime    time.Time
	LastModBy      int64
}

// Status values for Person
const (
	NORMAL     = int64(0)
	OPTOUT     = int64(1)
	BOUNCED    = int64(2)
	COMPLAINT  = int64(3)
	SUPPRESSED = int64(4)
)

// EGroup is a structure of all attributes of a EGroup to which Persons can belong
type EGroup struct {
	GID              int64
	GroupName        string
	GroupDescription string
	DtStart          time.Time // start time of last scrape
	DtStop           time.Time // stop time of last scrape
	LastModTime      time.Time // last time this record was updated
	LastModBy        int64
}

// PGroup is used to indicate that a person is a member of a group. A person can
// have multiple PGroups
type PGroup struct {
	PID         int64
	GID         int64
	LastModTime time.Time
	LastModBy   int64
}

// DataUpdate is used to describe a run where the data is updated.
// Every run generates this. So we have a history of updates
type DataUpdate struct {
	DUID        int64
	GID         int64     // group being scraped
	DtStart     time.Time // start time of a scrape
	DtStop      time.Time // stop time of a scrape
	LastModTime time.Time
	LastModBy   int64
}

// Query is a struct that represents a database query
// The QueryJSON is in a format that can be translated
// into SQL and used as a database query
type Query struct {
	QID         int64
	QueryName   string
	QueryDescr  string
	QueryJSON   string
	LastModTime time.Time
	LastModBy   int64
}

// DB is a struct of context data available for all the DB routines
var DB struct {
	Prepstmt PrepSQL
	Db       *sql.DB
	DBFields map[string]string
	Zone     *time.Location // what timezone should the server use?
}

// InitDB is the call to initialize database context set elsewhere
func InitDB(db *sql.DB) {
	var err error
	DB.Db = db
	DB.DBFields = map[string]string{}
	DB.Zone, err = time.LoadLocation(MojoDBConfig.Timezone)
	if err != nil {
		util.Console("Error loading timezone %s : %s\n", MojoDBConfig.Timezone, err.Error())
		util.Ulog("Error loading timezone %s : %s", MojoDBConfig.Timezone, err.Error())
	}
}

// GetJoinSetCount returns the number of database rows in the supplied table with the supplied where clause
func GetJoinSetCount(q string) (int64, error) {
	var count int64
	err := DB.Db.QueryRow(q).Scan(&count)
	return count, err
}

// GetRowCount returns the number of database rows in the supplied table with the supplied where clause
func GetRowCount(table, joins, where string) (int64, error) {
	if len(where) > 0 {
		where = " WHERE " + where
	}
	return GetRowCountRaw(table, joins, where)
}

// GetRowCountRaw returns the number of database rows in the supplied table with
// the supplied where clause. The where clause can be empty.
//
// INPUTS
//    table - table are we querying
//    joins - any join info, can be nil or an empty string
//    where - the where clause, can be nil or an empty string
//-----------------------------------------------------------------------------
func GetRowCountRaw(table, joins, where string) (int64, error) {
	count := int64(0)
	var err error
	s := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if len(joins) > 0 {
		s += " " + joins
	}
	if len(where) > 0 {
		s += " " + where
	}
	util.Console("\n\nGetRowCountRaw: QUERY = %s\n", s)
	de := DB.Db.QueryRow(s).Scan(&count)
	if de != nil {
		err = fmt.Errorf("GetRowCountRaw: query=\"%s\"    err = %s", s, de.Error())
	}
	return count, err
}

// GetQueryRowCount returns the number of rows in the solution set for the supplied named query
func GetQueryRowCount(qname string) (int64, error) {
	util.Console("entered GetQueryRowCount:  searching for %s\n", qname)
	c := int64(0)
	q, err := GetQueryByName(qname)
	if err != nil {
		return c, err
	}
	// THIS IS A HACK!!  Need to revamp when we get the real sql code in place
	// here's what a query looks like now:
	//		select People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=2 WHERE People.Status=0
	util.Console("Successfully read query: %s\n", q.QueryName)
	i := strings.Index(q.QueryJSON, "FROM")
	if i < 0 {
		return c, fmt.Errorf("could not find FROM in query")
	}
	s := "SELECT COUNT(People.PID) " + q.QueryJSON[i:]
	err = DB.Db.QueryRow(s).Scan(&c)
	return c, err
}

//=================================================
//                    PEOPLE
//=================================================
func readPerson(row *sql.Row) (Person, error) {
	var a Person
	err := row.Scan(&a.PID, &a.FirstName, &a.MiddleName, &a.LastName, &a.PreferredName, &a.JobTitle, &a.OfficePhone, &a.OfficeFax, &a.Email1, &a.Email2, &a.MailAddress, &a.MailAddress2, &a.MailCity, &a.MailState, &a.MailPostalCode, &a.MailCountry, &a.RoomNumber, &a.MailStop, &a.Status, &a.OptOutDate, &a.LastModTime, &a.LastModBy)
	return a, err
}

// GetPerson reads a Person the structure for the supplied id
func GetPerson(id int64) (Person, error) {
	return readPerson(DB.Prepstmt.GetPerson.QueryRow(id))
}

// GetPersonByEmail reads a Person the structure for the supplied email addr
func GetPersonByEmail(s1 string) (Person, error) {
	return readPerson(DB.Prepstmt.GetPersonByEmail.QueryRow(s1))
}

// ReadPersonFromRows uses the supplied sql.Rows struct to read a Person record
func ReadPersonFromRows(rows *sql.Rows) (Person, error) {
	var a Person
	err := rows.Scan(&a.PID, &a.FirstName, &a.MiddleName, &a.LastName, &a.PreferredName, &a.JobTitle, &a.OfficePhone, &a.OfficeFax, &a.Email1, &a.Email2, &a.MailAddress, &a.MailAddress2, &a.MailCity, &a.MailState, &a.MailPostalCode, &a.MailCountry, &a.RoomNumber, &a.MailStop, &a.Status, &a.OptOutDate, &a.LastModTime, &a.LastModBy)
	return a, err
}

// ReadPeople uses the supplied sql.Rows struct to read Person records into a slice.
// It returns the slice and any error encountered
func ReadPeople(rows *sql.Rows, err error) ([]Person, error) {
	var t []Person
	defer rows.Close()
	if err != nil {
		return t, err
	}
	for i := 0; rows.Next(); i++ {
		var a Person
		err = rows.Scan(&a.PID, &a.FirstName, &a.MiddleName, &a.LastName, &a.PreferredName, &a.JobTitle, &a.OfficePhone, &a.OfficeFax, &a.Email1, &a.Email2, &a.MailAddress, &a.MailAddress2, &a.MailCity, &a.MailState, &a.MailPostalCode, &a.MailCountry, &a.RoomNumber, &a.MailStop, &a.Status, &a.OptOutDate, &a.LastModTime, &a.LastModBy)
		if err != nil {
			return t, err
		}
		t = append(t, a)
	}
	return t, err
}

// GetPersonByName reads a Person the structure for the supplied id. Since there
// may be multiple people with the same name, an array of matches is returned.
func GetPersonByName(f, m, l string) ([]Person, error) {
	return ReadPeople(DB.Prepstmt.GetPersonByName.Query(f, m, l))
}

// GetPersonByRecordFieldMatching reads a Person the structure for the supplied id. Since there
// may be multiple people with the same name, an array of matches is returned.
func GetPersonByRecordFieldMatching(f, m, l, o, e string) ([]Person, error) {
	return ReadPeople(DB.Prepstmt.GetPersonByRecordFieldMatching.Query(f, f, m, l, o, e))
}

// UpdatePerson updates the existing database record for a
func UpdatePerson(a *Person) error {
	_, err := DB.Prepstmt.UpdatePerson.Exec(a.FirstName, a.MiddleName, a.LastName, a.PreferredName, a.JobTitle, a.OfficePhone, a.OfficeFax, a.Email1, a.Email2, a.MailAddress, a.MailAddress2, a.MailCity, a.MailState, a.MailPostalCode, a.MailCountry, a.RoomNumber, a.MailStop, a.Status, a.OptOutDate, a.LastModBy, a.PID)
	return err
}

// InsertPerson writes a new Person record to the database
func InsertPerson(a *Person) error {
	res, err := DB.Prepstmt.InsertPerson.Exec(a.FirstName, a.MiddleName, a.LastName, a.PreferredName, a.JobTitle, a.OfficePhone, a.OfficeFax, a.Email1, a.Email2, a.MailAddress, a.MailAddress2, a.MailCity, a.MailState, a.MailPostalCode, a.MailCountry, a.RoomNumber, a.MailStop, a.Status, a.OptOutDate, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			a.PID = int64(id)
		}
	} else {
		util.Ulog("InsertPerson: error inserting Person:  %v\n", err)
		util.Ulog("Person = %#v\n", *a)
	}
	return err
}

// DeletePerson deletes Person records with the supplied id
func DeletePerson(id int64) error {
	_, err := DB.Prepstmt.DeletePerson.Exec(id)
	if err != nil {
		util.Ulog("Error deleting Person for id = %d, error: %v\n", id, err)
	}
	return err
}

//=================================================
//                    EGROUP
//=================================================

// ReadGroup reads a row from the database EGroup table based on the supplied row
func ReadGroup(row *sql.Row) (EGroup, error) {
	var a EGroup
	err := row.Scan(&a.GID, &a.GroupName, &a.GroupDescription, &a.DtStart, &a.DtStop, &a.LastModTime, &a.LastModBy)
	return a, err
}

// GetGroup reads a EGroup the structure for the supplied id
func GetGroup(id int64) (EGroup, error) {
	return ReadGroup(DB.Prepstmt.GetGroup.QueryRow(id))
}

// GetGroupByName reads a EGroup the structure for the supplied group name
func GetGroupByName(s string) (EGroup, error) {
	return ReadGroup(DB.Prepstmt.GetGroupByName.QueryRow(s))
}

// GetGroupTypedown returns the values needed for typedown controls:
// input:   ctx
//            s - string or substring to search for
//        limit - return no more than this many matches
// return a slice of Groups and an error.
func GetGroupTypedown(ctx context.Context, s string, limit int) ([]EGroup, error) {
	var err error
	var m []EGroup
	var rows *sql.Rows

	s = "%" + s + "%"
	if rows, err = DB.Prepstmt.GetGroupTypedown.Query(s); err != nil {
		return m, err
	}
	defer rows.Close()

	for rows.Next() {
		var t EGroup
		if err = rows.Scan(&t.GID, &t.GroupName); err != nil {
			return m, err
		}
		m = append(m, t)
	}

	return m, rows.Err()
}

// ReadGroups reads one row from the supplied rows struct.
func ReadGroups(rows *sql.Rows) (EGroup, error) {
	var a EGroup
	err := rows.Scan(&a.GID, &a.GroupName, &a.GroupDescription, &a.DtStart, &a.DtStop, &a.LastModTime, &a.LastModBy)
	return a, err
}

// InsertGroup writes a new EGroup record to the database
func InsertGroup(a *EGroup) error {
	res, err := DB.Prepstmt.InsertGroup.Exec(a.GroupName, a.GroupDescription, a.DtStart, a.DtStop, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			a.GID = int64(id)
		}
	} else {
		util.Ulog("InsertGroup: error inserting EGroup:  %v\n", err)
		util.Ulog("EGroup = %#v\n", *a)
	}
	return err
}

// UpdateGroup updates the existing database record for a
func UpdateGroup(a *EGroup) error {
	_, err := DB.Prepstmt.UpdateGroup.Exec(a.GroupName, a.GroupDescription, a.DtStart, a.DtStop, a.LastModBy, a.GID)
	return err
}

// DeleteGroup deletes Group records with the supplied id
func DeleteGroup(id int64) error {
	_, err := DB.Prepstmt.DeleteGroup.Exec(id)
	if err != nil {
		util.Ulog("Error deleting Group for id = %d, error: %v\n", id, err)
	}
	return err
}

//=================================================
//                    PGROUP
//=================================================

func readPGroup(row *sql.Row) (PGroup, error) {
	var a PGroup
	err := row.Scan(&a.PID, &a.GID, &a.LastModTime, &a.LastModBy)
	return a, err
}

// GetPGroup reads a PGroup the structure for the supplied id
func GetPGroup(pid, gid int64) (PGroup, error) {
	return readPGroup(DB.Prepstmt.GetPGroup.QueryRow(pid, gid))
}

// InsertPGroup inserts a new PGroup record into the database
func InsertPGroup(a *PGroup) error {
	_, err := DB.Prepstmt.InsertPGroup.Exec(a.PID, a.GID, a.LastModBy)
	if nil != err {
		util.Ulog("InsertGroup: error inserting EGroup:  %v\n", err)
		util.Ulog("PGroup = %#v\n", *a)
	}
	return err
}

// DeletePGroup deletes a PGroup record from the database
func DeletePGroup(pid, gid int64) error {
	_, err := DB.Prepstmt.DeletePGroup.Exec(pid, gid)
	if err != nil {
		util.Ulog("Error deleting PGroup for pid = %d, gid = %d, error: %v\n", pid, gid, err)
	}
	return err
}

//=================================================
//               DATA UPDATE
//=================================================

func readDataUpdate(row *sql.Row) (DataUpdate, error) {
	var a DataUpdate
	err := row.Scan(&a.DUID, &a.DtStart, &a.DtStop, &a.LastModTime, &a.LastModBy)
	return a, err
}

// GetDataUpdate reads a DataUpdate the structure for the supplied id
func GetDataUpdate(id int64) (DataUpdate, error) {
	return readDataUpdate(DB.Prepstmt.GetDataUpdate.QueryRow(id))
}

// readDataUpdates uses the supplied sql.Rows struct to read DataUpdate records into a slice.
// It returns the slice and any error encountered
func readDataUpdates(rows *sql.Rows, err error) ([]DataUpdate, error) {
	var t []DataUpdate
	defer rows.Close()
	if err != nil {
		return t, err
	}
	for i := 0; rows.Next(); i++ {
		var a DataUpdate
		err = rows.Scan(&a.DUID, &a.DtStart, &a.DtStop, &a.LastModTime, &a.LastModBy)
		if err != nil {
			return t, err
		}
		t = append(t, a)
	}
	return t, err
}

// GetDataUpdateByGroup reads a DataUpdate the structure for the supplied id
func GetDataUpdateByGroup(id int64) ([]DataUpdate, error) {
	return readDataUpdates(DB.Prepstmt.GetDataUpdateByGroup.Query(id))
}

// InsertDataUpdate inserts a new DataUpdate record into the database
func InsertDataUpdate(a *DataUpdate) error {
	res, err := DB.Prepstmt.InsertDataUpdate.Exec(a.GID, a.DtStart, a.DtStop, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			a.DUID = int64(id)
		}
	} else {
		util.Ulog("InsertDataUpdate: error inserting DataUpdate:  %v\n", err)
		util.Ulog("DataUpdate = %#v\n", *a)
	}
	return err
}

// UpdateDataUpdate updates the existing DataUpdate record
func UpdateDataUpdate(a *DataUpdate) error {
	_, err := DB.Prepstmt.UpdateDataUpdate.Exec(a.GID, a.DtStart, a.DtStop, a.LastModBy, a.DUID)
	return err
}

//=================================================
//               QUERY
//=================================================

// ReadQuery reads a query record based on the supplied row
func ReadQuery(row *sql.Row) (Query, error) {
	var a Query
	err := row.Scan(&a.QID, &a.QueryName, &a.QueryDescr, &a.QueryJSON, &a.LastModTime, &a.LastModBy)
	return a, err
}

// ReadQueries reads the next query record based on the supplied rows
func ReadQueries(rows *sql.Rows) (Query, error) {
	var a Query
	err := rows.Scan(&a.QID, &a.QueryName, &a.QueryDescr, &a.QueryJSON, &a.LastModTime, &a.LastModBy)
	return a, err
}

// GetQuery reads a Query the structure for the supplied id
func GetQuery(id int64) (Query, error) {
	return ReadQuery(DB.Prepstmt.GetQuery.QueryRow(id))
}

// GetQueryByName reads a Query the structure for the supplied id
func GetQueryByName(s string) (Query, error) {
	return ReadQuery(DB.Prepstmt.GetQueryByName.QueryRow(s))
}

// InsertQuery inserts a new Query record into the database
func InsertQuery(a *Query) error {
	res, err := DB.Prepstmt.InsertQuery.Exec(a.QueryName, a.QueryDescr, a.QueryJSON, a.LastModBy)
	if nil == err {
		id, err := res.LastInsertId()
		if err == nil {
			a.QID = int64(id)
		}
	} else {
		util.Ulog("InsertQuery: error inserting Query:  %v\n", err)
		util.Ulog("Query = %#v\n", *a)
	}
	return err
}

// UpdateQuery updates the existing Query record
func UpdateQuery(a *Query) error {
	_, err := DB.Prepstmt.UpdateQuery.Exec(a.QueryName, a.QueryDescr, a.QueryJSON, a.LastModBy, a.QID)
	if err != nil {
		util.Ulog("InsertQuery: error updating Query:  %v\n", err)
	}
	return err
}

// DeleteQuery deletes a Query record from the database
func DeleteQuery(id int64) error {
	_, err := DB.Prepstmt.DeleteQuery.Exec(id)
	if err != nil {
		util.Ulog("Error deleting Query for id = %d, error: %v\n", id, err)
	}
	return err
}
