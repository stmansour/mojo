package main

import (
	"encoding/json"
	"fmt"
	"mojo/db"
	"mojo/util"
	"net/http"
	"strings"
	"time"
)

// QueryGrid contains the data from Query that is targeted to the UI Grid that displays
// a list of Query structs
type QueryGrid struct {
	Recid       int64 `json:"recid"`
	QID         int64
	QueryName   string
	QueryDescr  string
	QueryJSON   string
	LastModTime time.Time
	LastModBy   int64
}

// QuerySearchResponse is a response string to the search request for Query records
type QuerySearchResponse struct {
	Status  string      `json:"status"`
	Total   int64       `json:"total"`
	Records []QueryGrid `json:"records"`
}

// QueryGridSave is the input data format for a Save command
type QueryGridSave struct {
	Status   string      `json:"status"`
	Recid    int64       `json:"recid"`
	FormName string      `json:"name"`
	Record   QueryGrid   `json:"record"`
	Changes  []QueryGrid `json:"changes"`
}

// QueryGetResponse is the response to a GetQuery request
type QueryGetResponse struct {
	Status string    `json:"status"`
	Record QueryGrid `json:"record"`
}

// QueryStats is a structure some interesting statistics for the Query table
type QueryStats struct {
	MemberCount     int64
	LastScrapeStart string
	LastScrapeStop  string
}

// QueryStatResponse is the response to a Query stats request
type QueryStatResponse struct {
	Status string     `json:"status"`
	Record QueryStats `json:"record"`
}

// SvcHandlerQuery formats a complete data record for an assessment for use with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
// The server command can be:
//      get
//      save
//      delete
//-----------------------------------------------------------------------------------
func SvcHandlerQuery(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	fmt.Printf("Entered SvcHandlerQuery\n")

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchHandlerQueries(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcGridErrorReturn(w, fmt.Errorf("QueryID is required but was not specified"))
				return
			}
			getQuery(w, r, d)
		}
		break
	case "save":
		saveQuery(w, r, d)
		break
	case "delete":
		deleteQuery(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcGridErrorReturn(w, err)
		return
	}
}

// SvcQueryResultsCount returns the number of people in the database
// wsdoc {
//  @Title  QueryResults Count
//	@URL /v1/qrescount/QID
//  @Method  GET
//	@Synopsis Get the count of people in the database
//  @Descr  Returns a count of all people in the database. If GID
//  @Descr  is provided it returns the count of people in group GID.
//	@Input WebGridSearchRequest
//  @Response CountResponse
// wsdoc }
func SvcQueryResultsCount(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcQueryResultsCount"
	fmt.Printf("Entered %s\n", funcname)
	var (
		g   CountResponse
		err error
	)
	if d.ID < 0 {
		e := fmt.Errorf("%s: Required parameter QID was not provided", funcname)
		SvcGridErrorReturn(w, e)
		return
	}
	q, err := db.GetQuery(d.ID)
	if err != nil {
		e := fmt.Errorf("%s: Error getting Query: %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	// ex:  SELECT People.*,x,y FROM People INNER JOIN PGroup ON PGroup.PID=People.PID && PGroup.GID=%d WHERE People.Status=0",d.ID)
	// TODO: this needs to be redone in when we generalize the query
	i := strings.Index(strings.ToLower(q.QueryJSON), "from")
	if i < 0 {
		e := fmt.Errorf("%s: Could not find 'FROM' in Query: %s", funcname, q.QueryJSON)
		SvcGridErrorReturn(w, e)
		return
	}
	s := q.QueryJSON[i+5:] // moves us up to the table name
	i = strings.Index(s, " ")
	if i < 0 {
		e := fmt.Errorf("%s: No space after table name in Query: %s", funcname, q.QueryJSON)
		SvcGridErrorReturn(w, e)
		return
	}
	qw := s[i+1:] // this is the part of the query we need...

	g.Record.Count, err = GetRowCountRaw("People", qw)
	if err != nil {
		fmt.Printf("Error from GetRowCount: %s\n", err.Error())
		SvcGridErrorReturn(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// SvcSearchHandlerQueries generates a report of all Queries defined business d.BID
// wsdoc {
//  @Title  Search Queries
//	@URL /v1/Queries/[:QID]
//  @Method  POST
//	@Synopsis Search Queries
//  @Descr  Search all Query and return those that match the Search Logic.
//  @Descr  The search criteria includes start and stop dates of interest.
//	@Input WebGridSearchRequest
//  @Response QuerySearchResponse
// wsdoc }
func SvcSearchHandlerQueries(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchHandlerQueries"
	fmt.Printf("Entered %s\n", funcname)
	var (
		g   QuerySearchResponse
		err error
	)

	order := "QueryName ASC"                                           // default ORDER
	q := fmt.Sprintf("SELECT %s FROM Query ", db.DB.DBFields["Query"]) // the fields we want
	qw := fmt.Sprintf("")                                              // don't need WHERE clause on this query
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
	fmt.Printf("rowcount query conditions: %s\ndb query = %s\n", qw, q)

	g.Total, err = GetRowCount("Query", qw)
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
		var q QueryGrid
		p, err := db.ReadQueries(rows)
		if err != nil {
			fmt.Printf("%s.  Error reading Query: %s\n", funcname, err.Error())
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

// deleteQuery deletes a payment type from the database
// wsdoc {
//  @Title  Delete Query
//	@URL /v1/dep/:BUI/:RAID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a Query.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
func deleteQuery(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteQuery"
	fmt.Printf("Entered %s\n", funcname)
	fmt.Printf("record data = %s\n", d.data)
	var del WebGridDelete
	if err := json.Unmarshal([]byte(d.data), &del); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	for i := 0; i < len(del.Selected); i++ {
		if err := db.DeleteQuery(del.Selected[i]); err != nil {
			SvcGridErrorReturn(w, err)
			return
		}
	}
	SvcWriteSuccessResponse(w)
}

// GetQuery returns the requested assessment
// wsdoc {
//  @Title  Save Query
//	@URL /v1/dep/:BUI/:PID
//  @Method  GET
//	@Synopsis Update the information on a Query with the supplied data
//  @Description  This service updates Query :PID with the information supplied. All fields must be supplied.
//	@Input QueryGridSave
//  @Response SvcStatusResponse
// wsdoc }
func saveQuery(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveQuery"
	fmt.Printf("Entered %s\n", funcname)
	fmt.Printf("record data = %s\n", d.data)

	var foo QueryGridSave
	data := []byte(d.data)
	err := json.Unmarshal(data, &foo)

	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	if len(foo.Changes) == 0 { // This is a new record
		var a db.Query
		util.MigrateStructVals(&foo.Record, &a) // the variables that don't need special handling
		fmt.Printf("a = %#v\n", a)
		fmt.Printf(">>>> NEW PAYMENT TYPE IS BEING ADDED\n")
		err = db.InsertQuery(&a)
		if err != nil {
			e := fmt.Errorf("%s: Error saving Query: %s", funcname, err.Error())
			SvcGridErrorReturn(w, e)
			return
		}
	} else { // update existing or add new record(s)
		fmt.Printf("Uh oh - we have not yet implemented this!!!\n")
		fmt.Fprintf(w, "Have not implemented this function")
		// if err = JSONchangeParseUtil(d.data, QueryUpdate, d); err != nil {
		// 	SvcGridErrorReturn(w, err)
		// 	return
		// }
	}
	SvcWriteSuccessResponse(w)
}

// QueryUpdate unmarshals the supplied string. If Recid > 0 it updates the
// Query record using Recid as the PID.  If Recid == 0, then it inserts a
// new Query record.
func QueryUpdate(s string, d *ServiceData) error {
	var err error
	b := []byte(s)
	var rec QueryGrid
	if err = json.Unmarshal(b, &rec); err != nil { // first parse to determine the record ID we need to load
		return err
	}
	if rec.Recid > 0 { // is this an update?
		pt, err := db.GetQuery(rec.Recid) // now load that record...
		if err != nil {
			return err
		}
		if err = json.Unmarshal(b, &pt); err != nil { // merge in the changes...
			return err
		}
		return db.UpdateQuery(&pt) // and save the result
	}
	// no, it is a new table entry that has not been saved...
	var a db.Query
	if err := json.Unmarshal(b, &a); err != nil { // merge in the changes...
		return err
	}
	fmt.Printf("a = %#v\n", a)
	fmt.Printf(">>>> NEW Query IS BEING ADDED\n")
	err = db.InsertQuery(&a)
	return err
}

// GetQuery returns the requested assessment
// wsdoc {
//  @Title  Get Query
//	@URL /v1/dep/:BUI/:PID
//  @Method  GET
//	@Synopsis Get information on a Query
//  @Description  Return all fields for assessment :PID
//	@Input WebGridSearchRequest
//  @Response QueryGetResponse
// wsdoc }
func getQuery(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getQuery"
	fmt.Printf("entered %s\n", funcname)
	var g QueryGetResponse
	a, err := db.GetQuery(d.ID)
	if err != nil {
		SvcGridErrorReturn(w, err)
		return
	}
	if a.QID > 0 {
		var gg QueryGrid
		util.MigrateStructVals(&a, &gg)
		g.Record = gg
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
