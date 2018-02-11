package main

import (
	"encoding/json"
	"fmt"
	"mojo/db"
	"mojo/util"
	"net/http"
	"time"
)

// PersonGrid contains the data from Person that is targeted to the UI Grid that displays
// a list of Person structs
type PersonGrid struct {
	Recid          int64 `json:"recid"`
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
	OptOutDate     util.JSONDate /*Time*/
	LastModTime    time.Time
	LastModBy      int64
}

// PersonSearchResponse is a response string to the search request for Person records
type PersonSearchResponse struct {
	Status  string       `json:"status"`
	Total   int64        `json:"total"`
	Records []PersonGrid `json:"records"`
}

// PersonGridSave is the input data format for a Save command
type PersonGridSave struct {
	Status   string       `json:"status"`
	Recid    int64        `json:"recid"`
	FormName string       `json:"name"`
	Record   PersonGrid   `json:"record"`
	Changes  []PersonGrid `json:"changes"`
}

// RecordCount is a structure with the count of records
// for some particular table.
type RecordCount struct {
	Recid int64 `json:"recid"`
	Count int64
}

// CountResponse is the response to a PersonCount request
type CountResponse struct {
	Status string      `json:"status"`
	Record RecordCount `json:"record"`
}

// PeopleStats is a structure some interesting statistics for the People table
type PeopleStats struct {
	Count      int64
	OptOut     int64
	Bounced    int64
	Complaint  int64
	Suppressed int64
}

// PeopleStatResponse is the response to a PersonCount request
type PeopleStatResponse struct {
	Status string      `json:"status"`
	Record PeopleStats `json:"record"`
}

// PersonGetResponse is the response to a GetPerson request
type PersonGetResponse struct {
	Status string     `json:"status"`
	Record PersonGrid `json:"record"`
}

// SvcHandlerPerson formats a complete data record for an assessment for use with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
// The server command can be:
//      get
//      save
//      delete
//-----------------------------------------------------------------------------------
func SvcHandlerPerson(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerPerson\n")

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchHandlerPeople(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcGridErrorReturn(w, fmt.Errorf("PersonID is required but was not specified"))
				return
			}
			getPerson(w, r, d)
		}
		break
	case "save":
		savePerson(w, r, d)
		break
	case "delete":
		deletePerson(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcGridErrorReturn(w, err)
		return
	}
}

// SvcPeopleCount returns the number of people in the database
// wsdoc {
//  @Title  People Count
//	@URL /v1/peoplecount/[GID]
//  @Method  POST GET
//	@Synopsis Get the count of people in the database
//  @Descr  Returns a count of all people in the database. If GID
//  @Descr  is provided it returns the count of people in group GID.
//	@Input WebGridSearchRequest
//  @Response CountResponse
// wsdoc }
//-----------------------------------------------------------------------------
func SvcPeopleCount(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcPeopleCount"
	util.Console("Entered %s\n", funcname)
	var (
		g   CountResponse
		err error
	)

	g.Record.Count, err = db.GetRowCount("People", "")
	if err != nil {
		util.Console("Error from db.GetRowCount: %s\n", err.Error())
		SvcGridErrorReturn(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// SvcPeopleStats returns statistics on the people table in the database
// wsdoc {
//  @Title  People Stats
//	@URL /v1/peoplestats/[:GID]
//  @Method  POST GET
//	@Synopsis Get statistics on the people table.
//  @Descr  Returns a count all the people, how many have opted out
//  @Descr  and how many email addresses have bounced.
//	@Input WebGridSearchRequest
//  @Response CountResponse
// wsdoc }
//-----------------------------------------------------------------------------
func SvcPeopleStats(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcPeopleStats"
	util.Console("Entered %s\n", funcname)
	var (
		g   PeopleStatResponse
		err error
	)

	s := ""
	if d.ID > 0 {
		s = fmt.Sprintf("WHERE GID=%d AND ", d.ID)
	}

	var m = []struct {
		Where string
		Count *int64
	}{
		{s, &g.Record.Count},
		{s + "Status=1", &g.Record.OptOut},
		{s + "Status=2", &g.Record.Bounced},
		{s + "Status=3", &g.Record.Complaint},
		{s + "Status=4", &g.Record.Suppressed},
	}

	for i := 0; i < len(m); i++ {
		*m[i].Count, err = db.GetRowCount("People", m[i].Where)
		if err != nil {
			util.Console("Error from db.GetRowCount: i = %d, err: %s\n", i, err.Error())
			SvcGridErrorReturn(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// SvcSearchHandlerPeople generates a report of all People defined business d.BID
// wsdoc {
//  @Title  Search People
//	@URL /v1/people/[:GID]
//  @Method  POST
//	@Synopsis Search People
//  @Descr  Search all Person and return those that match the Search Logic.
//  @Descr  The search criteria includes start and stop dates of interest.
//	@Input WebGridSearchRequest
//  @Response PersonSearchResponse
// wsdoc }
//-----------------------------------------------------------------------------
func SvcSearchHandlerPeople(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchHandlerPeople"
	util.Console("Entered %s\n", funcname)
	var (
		g   PersonSearchResponse
		err error
	)

	order := "PID ASC"                                                   // default ORDER
	q := fmt.Sprintf("SELECT %s FROM People ", db.DB.DBFields["People"]) // the fields we want
	qw := ""
	if len(d.wsSearchReq.Search) > 0 {
		v := d.wsSearchReq.Search[0].Value
		qw = fmt.Sprintf("FirstName LIKE \"%%%s%%\" OR LastName LIKE \"%%%s%%\" OR Email1 LIKE \"%%%s%%\"", v, v, v)
	}
	if len(qw) > 0 {
		q += "WHERE " + qw
	}
	q += " ORDER BY "
	if len(d.wsSearchReq.Sort) > 0 {
		for i := 0; i < len(d.wsSearchReq.Sort); i++ {
			if i > 0 {
				q += ","
			}
			q += d.wsSearchReq.Sort[i].Field + " " + d.wsSearchReq.Sort[i].Direction
		}
	} else {
		q += order
	}

	// now set up the offset and limit
	q += fmt.Sprintf(" LIMIT %d OFFSET %d", d.wsSearchReq.Limit, d.wsSearchReq.Offset)
	util.Console("rowcount query conditions: %s\ndb query = %s\n", qw, q)

	g.Total, err = db.GetRowCount("People", qw)
	if err != nil {
		util.Console("Error from db.GetRowCount: %s\n", err.Error())
		SvcGridErrorReturn(w, err)
		return
	}
	rows, err := db.DB.Db.Query(q)
	if err != nil {
		util.Console("Error from DB Query: %s\n", err.Error())
		SvcGridErrorReturn(w, err)
		return
	}
	defer rows.Close()

	i := int64(d.wsSearchReq.Offset)
	count := 0
	for rows.Next() {
		var q PersonGrid
		p, err := db.ReadPersonFromRows(rows)
		if err != nil {
			util.Console("%s.  Error reading Person: %s\n", funcname, err.Error())
		}
		util.MigrateStructVals(&p, &q)
		q.Recid = p.PID
		g.Records = append(g.Records, q)
		count++ // update the count only after adding the record
		if count >= d.wsSearchReq.Limit {
			break // if we've added the max number requested, then exit
		}
		i++
	}
	util.Console("g.Total = %d\n", g.Total)
	util.ErrCheck(rows.Err())
	w.Header().Set("Content-Type", "application/json")
	g.Status = "success"
	SvcWriteResponse(&g, w)

}

// deletePerson deletes a payment type from the database
// wsdoc {
//  @Title  Delete Person
//	@URL /v1/person/PID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a Person.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func deletePerson(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deletePerson"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	var del WebGridDelete
	if err := json.Unmarshal([]byte(d.data), &del); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	// Cmd      string  `json:"cmd"`
	// Selected []int64 `json:"selected"`
	// Limit    int     `json:"limit"`
	// Offset   int     `json:"offset"`

	for i := 0; i < len(del.Selected); i++ {
		if err := db.DeletePerson(del.Selected[i]); err != nil {
			SvcGridErrorReturn(w, err)
			return
		}
	}
	SvcWriteSuccessResponse(w)
}

// SavePerson returns the requested assessment
// wsdoc {
//  @Title  Save Person
//	@URL /v1/persone/PID
//  @Method  GET
//	@Synopsis Update the information on a Person with the supplied data, create if necessary.
//  @Description  This service creates a person if PID == 0 or updates a Person if PID > 0 with
//  @Description  the information supplied. All fields must be supplied.
//	@Input PersonGridSave
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func savePerson(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "savePerson"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)

	var foo PersonGridSave
	data := []byte(d.data)
	err := json.Unmarshal(data, &foo)

	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	if foo.Record.PID == 0 { // This is a new record
		var a db.Person
		util.MigrateStructVals(&foo.Record, &a) // the variables that don't need special handling
		util.Console("a = %#v\n", a)
		util.Console(">>>> NEW PERSON IS BEING ADDED\n")
		err = db.InsertPerson(&a)
		if err != nil {
			e := fmt.Errorf("%s: Error saving Person: %s", funcname, err.Error())
			SvcGridErrorReturn(w, e)
			return
		}
	} else { // update existing or add new record(s)
		err = PersonUpdate(&foo.Record, d)
		if err != nil {
			SvcGridErrorReturn(w, err)
			return
		}
	}
	SvcWriteSuccessResponse(w)
}

// PersonUpdate updates the supplied person in the database with the supplied
// info. It only allows certain fields to be updated.
//-----------------------------------------------------------------------------
func PersonUpdate(p *PersonGrid, d *ServiceData) error {
	util.Console("entered PersonUpdate, p.Email2 = %s\n", p.Email2)
	var err error
	pt, err := db.GetPerson(p.PID) // now load that record...
	if err != nil {
		return err
	}

	if !util.ValidEmailAddress(p.Email1) {
		return fmt.Errorf("Invalid email address: %s", p.Email1)
	}

	if !util.ValidEmailAddress(p.Email2) {
		return fmt.Errorf("Invalid email address: %s", p.Email2)
	}

	pt.FirstName = p.FirstName
	pt.MiddleName = p.MiddleName
	pt.LastName = p.LastName
	pt.PreferredName = p.PreferredName
	pt.JobTitle = p.JobTitle
	pt.OfficePhone = p.OfficePhone
	pt.OfficeFax = p.OfficeFax
	pt.Email1 = p.Email1
	pt.Email2 = p.Email2
	pt.MailAddress = p.MailAddress
	pt.MailAddress2 = p.MailAddress2
	pt.MailCity = p.MailCity
	pt.MailState = p.MailState
	pt.MailPostalCode = p.MailPostalCode
	pt.MailCountry = p.MailCountry
	pt.RoomNumber = p.RoomNumber
	pt.MailStop = p.MailStop
	pt.Status = p.Status

	return db.UpdatePerson(&pt) // and save the result
}

// GetPerson returns the requested assessment
// wsdoc {
//  @Title  Get Person
//	@URL /v1/dep/:BUI/:PID
//  @Method  GET
//	@Synopsis Get information on a Person
//  @Description  Return all fields for assessment :PID
//	@Input WebGridSearchRequest
//  @Response PersonGetResponse
// wsdoc }
//-----------------------------------------------------------------------------
func getPerson(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getPerson"
	util.Console("entered %s\n", funcname)
	var g PersonGetResponse
	a, err := db.GetPerson(d.ID)
	if err != nil {
		SvcGridErrorReturn(w, err)
		return
	}
	if a.PID > 0 {
		var gg PersonGrid
		util.MigrateStructVals(&a, &gg)
		gg.Recid = gg.PID
		g.Record = gg
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
