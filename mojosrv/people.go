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

// PersonGetResponse is the response to a GetPerson request
type PersonGetResponse struct {
	Status string     `json:"status"`
	Record PersonGrid `json:"record"`
}

// GetRowCount returns the number of database rows in the supplied table with the supplied where clause
func GetRowCount(table, where string) (int64, error) {
	count := int64(0)
	var err error
	s := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, where)
	de := db.DB.Db.QueryRow(s).Scan(&count)
	if de != nil {
		err = fmt.Errorf("GetRowCount: query=\"%s\"    err = %s", s, de.Error())
	}
	return count, err
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
	fmt.Printf("Entered SvcHandlerPerson\n")

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

// SvcSearchHandlerPeople generates a report of all People defined business d.BID
// wsdoc {
//  @Title  Search People
//	@URL /v1/dep/:BUI
//  @Method  POST
//	@Synopsis Search People
//  @Descr  Search all Person and return those that match the Search Logic.
//  @Descr  The search criteria includes start and stop dates of interest.
//	@Input WebGridSearchRequest
//  @Response PersonSearchResponse
// wsdoc }
func SvcSearchHandlerPeople(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchHandlerPeople"
	fmt.Printf("Entered %s\n", funcname)
	var (
		g   PersonSearchResponse
		err error
	)

	order := "PID ASC"                                                   // default ORDER
	q := fmt.Sprintf("SELECT %s FROM People ", db.DB.DBFields["People"]) // the fields we want
	qw := fmt.Sprintf("PID>0")                                           // will probably change this at some point
	q += "WHERE " + qw + " ORDER BY "
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
	fmt.Printf("rowcount query conditions: %s\ndb query = %s\n", qw, q)

	g.Total, err = GetRowCount("People", qw)
	if err != nil {
		fmt.Printf("Error from GetRowCount: %s\n", err.Error())
		SvcGridErrorReturn(w, err)
		return
	}
	rows, err := db.DB.Db.Query(q)
	if err != nil {
		fmt.Printf("Error from DB Query: %s\n", err.Error())
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
			fmt.Printf("%s.  Error reading Person: %s\n", funcname, err.Error())
		}
		util.MigrateStructVals(&p, &q)
		g.Records = append(g.Records, q)
		count++ // update the count only after adding the record
		if count >= d.wsSearchReq.Limit {
			break // if we've added the max number requested, then exit
		}
		i++
	}
	fmt.Printf("g.Total = %d\n", g.Total)
	util.ErrCheck(rows.Err())
	w.Header().Set("Content-Type", "application/json")
	g.Status = "success"
	SvcWriteResponse(&g, w)

}

// deletePerson deletes a payment type from the database
// wsdoc {
//  @Title  Delete Person
//	@URL /v1/dep/:BUI/:RAID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a Person.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
func deletePerson(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deletePerson"
	fmt.Printf("Entered %s\n", funcname)
	fmt.Printf("record data = %s\n", d.data)
	var del WebGridDelete
	if err := json.Unmarshal([]byte(d.data), &del); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	for i := 0; i < len(del.Selected); i++ {
		if err := db.DeletePerson(del.Selected[i]); err != nil {
			SvcGridErrorReturn(w, err)
			return
		}
	}
	SvcWriteSuccessResponse(w)
}

// GetPerson returns the requested assessment
// wsdoc {
//  @Title  Save Person
//	@URL /v1/dep/:BUI/:PID
//  @Method  GET
//	@Synopsis Update the information on a Person with the supplied data
//  @Description  This service updates Person :PID with the information supplied. All fields must be supplied.
//	@Input PersonGridSave
//  @Response SvcStatusResponse
// wsdoc }
func savePerson(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "savePerson"
	fmt.Printf("Entered %s\n", funcname)
	fmt.Printf("record data = %s\n", d.data)

	var foo PersonGridSave
	data := []byte(d.data)
	err := json.Unmarshal(data, &foo)

	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	if len(foo.Changes) == 0 { // This is a new record
		var a db.Person
		util.MigrateStructVals(&foo.Record, &a) // the variables that don't need special handling
		fmt.Printf("a = %#v\n", a)
		fmt.Printf(">>>> NEW PAYMENT TYPE IS BEING ADDED\n")
		err = db.InsertPerson(&a)
		if err != nil {
			e := fmt.Errorf("%s: Error saving Person: %s", funcname, err.Error())
			SvcGridErrorReturn(w, e)
			return
		}
	} else { // update existing or add new record(s)
		fmt.Printf("Uh oh - we have not yet implemented this!!!\n")
		fmt.Fprintf(w, "Have not implemented this function")
		// if err = JSONchangeParseUtil(d.data, PersonUpdate, d); err != nil {
		// 	SvcGridErrorReturn(w, err)
		// 	return
		// }
	}
	SvcWriteSuccessResponse(w)
}

// PersonUpdate unmarshals the supplied string. If Recid > 0 it updates the
// Person record using Recid as the PID.  If Recid == 0, then it inserts a
// new Person record.
func PersonUpdate(s string, d *ServiceData) error {
	var err error
	b := []byte(s)
	var rec PersonGrid
	if err = json.Unmarshal(b, &rec); err != nil { // first parse to determine the record ID we need to load
		return err
	}
	if rec.Recid > 0 { // is this an update?
		pt, err := db.GetPerson(rec.Recid) // now load that record...
		if err != nil {
			return err
		}
		if err = json.Unmarshal(b, &pt); err != nil { // merge in the changes...
			return err
		}
		return db.UpdatePerson(&pt) // and save the result
	}
	// no, it is a new table entry that has not been saved...
	var a db.Person
	if err := json.Unmarshal(b, &a); err != nil { // merge in the changes...
		return err
	}
	fmt.Printf("a = %#v\n", a)
	fmt.Printf(">>>> NEW Person IS BEING ADDED\n")
	err = db.InsertPerson(&a)
	return err
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
func getPerson(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getPerson"
	fmt.Printf("entered %s\n", funcname)
	var g PersonGetResponse
	a, err := db.GetPerson(d.ID)
	if err != nil {
		SvcGridErrorReturn(w, err)
		return
	}
	if a.PID > 0 {
		var gg PersonGrid
		util.MigrateStructVals(&a, &gg)
		g.Record = gg
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
