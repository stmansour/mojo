package main

import (
	"encoding/json"
	"fmt"
	"mojo/db"
	"mojo/util"
	"net/http"
	"time"
)

// GroupGrid contains the data from Group that is targeted to the UI Grid that displays
// a list of Group structs
type GroupGrid struct {
	Recid            int64 `json:"recid"`
	GID              int64
	GroupName        string
	GroupDescription string
	LastModTime      time.Time
	LastModBy        int64
}

// GroupSave is a struct for saving a group
type GroupSave struct {
	Cmd    string    `json:"cmd"`
	Record GroupGrid `json:"record"`
}

// GroupSearchResponse is a response string to the search request for Group records
type GroupSearchResponse struct {
	Status  string      `json:"status"`
	Total   int64       `json:"total"`
	Records []GroupGrid `json:"records"`
}

// GroupGetResponse is the response to a GetGroup request
type GroupGetResponse struct {
	Status string    `json:"status"`
	Record GroupGrid `json:"record"`
}

// GroupStats is a structure some interesting statistics for the Group table
type GroupStats struct {
	GID              int64
	GroupName        string
	GroupDescription string
	MemberCount      int64
	MailToCount      int64
	OptOutCount      int64
	BouncedCount     int64
	ComplaintCount   int64
	SuppressedCount  int64
	LastScrapeStart  string
	LastScrapeStop   string
}

// GroupStatResponse is the response to a Group stats request
type GroupStatResponse struct {
	Status string     `json:"status"`
	Record GroupStats `json:"record"`
}

// SvcHandlerGroup formats a complete data record for an assessment for use with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
// The server command can be:
//      get
//      save
//      delete
//-----------------------------------------------------------------------------------
func SvcHandlerGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerGroup:  cmd = %s\n", d.wsSearchReq.Cmd)

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchHandlerGroups(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcGridErrorReturn(w, fmt.Errorf("GroupID is required but was not specified"))
				return
			}
			getGroup(w, r, d)
		}
	case "save":
		saveGroup(w, r, d)
	case "delete":
		deleteGroup(w, r, d)
	default:
		err := fmt.Errorf("unhandled command: %s", d.wsSearchReq.Cmd)
		SvcGridErrorReturn(w, err)
		return
	}
}

// SvcGroupsCount returns the number of people in the database
// wsdoc {
//  @Title  Groups Count
//	@URL /v1/groupcount/
//  @Method  POST GET
//	@Synopsis Get the count of people in the database
//  @Descr  Returns a count of all people in the database. If GID
//  @Descr  is provided it returns the count of people in grop GID
//	@Input WebGridSearchRequest
//  @Response GroupSearchResponse
// wsdoc }
func SvcGroupsCount(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchHandlerGroups"
	util.Console("Entered %s\n", funcname)
	var (
		g   CountResponse
		err error
	)

	g.Record.Count, err = db.GetRowCount("EGroup", "", "")
	if err != nil {
		util.Console("Error from db.GetRowCount: %s\n", err.Error())
		SvcGridErrorReturn(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// SvcSearchHandlerGroups gets a list of all groups matching the supplied
// search criteria.
//------------------------------------------------------------------------------
func SvcSearchHandlerGroups(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchHandlerGroups"
	util.Console("Entered %s\n", funcname)
	var (
		g   GroupSearchResponse
		err error
	)

	order := "GroupName ASC"                                             // default ORDER
	q := fmt.Sprintf("SELECT %s FROM EGroup ", db.DB.DBFields["EGroup"]) // the fields we want
	qw := ""                                           // don't need WHERE clause on this query
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

	g.Total, err = db.GetRowCount("EGroup", "", qw)
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
		var q GroupGrid
		p, err := db.ReadGroups(rows)
		if err != nil {
			util.Console("%s.  Error reading Group: %s\n", funcname, err.Error())
		}
		util.MigrateStructVals(&p, &q)
		q.Recid = q.GID
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

// deleteGroup deletes a payment type from the database
// wsdoc {
//  @Title  Delete Group
//	@URL /v1/dep/:BUI/:RAID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a Group.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
func deleteGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteGroup"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	var del WebGridDelete
	if err := json.Unmarshal([]byte(d.data), &del); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	for i := 0; i < len(del.Selected); i++ {
		if err := db.DeleteGroup(del.Selected[i]); err != nil {
			SvcGridErrorReturn(w, err)
			return
		}
	}
	SvcWriteSuccessResponse(w)
}

// saveGroup - saves the supplied group information, creating the group if needed.
//
//--------------------------------------------------------------------------------
func saveGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveGroup"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)

	var foo GroupSave
	var err error
	var grp db.EGroup
	data := []byte(d.data)
	err = json.Unmarshal(data, &foo)

	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}

	if foo.Record.GID == 0 {
		grp, err = db.GetGroupByName(foo.Record.GroupName)
	} else {
		grp, err = db.GetGroup(foo.Record.GID)
	}
	if nil != err {
		if util.IsSQLNoResultsError(err) {
			grp.GroupName = foo.Record.GroupName
			grp.DtStart = time.Now()
			grp.GroupDescription = foo.Record.GroupDescription
			if err = db.InsertGroup(&grp); err != nil {
				e := fmt.Errorf("error inserting group %s: %s", foo.Record.GroupName, err.Error())
				SvcGridErrorReturn(w, e)
				return
			}
		}
	} else {
		chg := false
		if grp.GroupName != foo.Record.GroupName {
			chg = true
			grp.GroupName = foo.Record.GroupName
		}
		if grp.GroupDescription != foo.Record.GroupDescription {
			chg = true
			grp.GroupDescription = foo.Record.GroupDescription
		}
		if chg {
			if err = db.UpdateGroup(&grp); err != nil {
				e := fmt.Errorf("error updating group %s: %s", foo.Record.GroupName, err.Error())
				SvcGridErrorReturn(w, e)
				return
			}
		}
	}

	SvcWriteSuccessResponseWithID(w, grp.GID)
}

// GroupUpdate unmarshals the supplied string. If Recid > 0 it updates the
// Group record using Recid as the PID.  If Recid == 0, then it inserts a
// new Group record.
func GroupUpdate(s string, d *ServiceData) error {
	var err error
	b := []byte(s)
	var rec GroupGrid
	if err = json.Unmarshal(b, &rec); err != nil { // first parse to determine the record ID we need to load
		return err
	}
	if rec.Recid > 0 { // is this an update?
		pt, err := db.GetGroup(rec.Recid) // now load that record...
		if err != nil {
			return err
		}
		if err = json.Unmarshal(b, &pt); err != nil { // merge in the changes...
			return err
		}
		return db.UpdateGroup(&pt) // and save the result
	}
	// no, it is a new table entry that has not been saved...
	var a db.EGroup
	if err := json.Unmarshal(b, &a); err != nil { // merge in the changes...
		return err
	}
	util.Console("a = %#v\n", a)
	util.Console(">>>> NEW Group IS BEING ADDED\n")
	err = db.InsertGroup(&a)
	return err
}

// GetGroupStats returns the requested assessment
// wsdoc {
//  @Title  Get Group Statistics
//	@URL /v1/groupstats/GID
//  @Method  GET
//	@Synopsis Get information and stats on a Group
//  @Description  Return all fields and solution set count for the supplied GID
//	@Input WebGridSearchRequest
//  @Response GroupGetResponse
// wsdoc }
func GetGroupStats(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const (
		DATETIMEINPFMT = "2006-01-02 15:04:00 MST"
	)
	funcname := "GetGroupStat"
	util.Console("entered %s.  Group id = %d\n", funcname, d.ID)
	var g GroupStatResponse
	a, err := db.GetGroup(d.ID)
	if err != nil {
		SvcGridErrorReturn(w, err)
		return
	}
	if a.GID > 0 {
		var gg GroupStats
		util.MigrateStructVals(&a, &gg)
		gg.LastScrapeStart = a.DtStart.Format(DATETIMEINPFMT)
		gg.LastScrapeStop = a.DtStop.Format(DATETIMEINPFMT)
		g.Record = gg
	}

	var gstat = []struct {
		q string
		r *int64
	}{
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d", r: &g.Record.MemberCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0", r: &g.Record.MailToCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=1", r: &g.Record.OptOutCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=2", r: &g.Record.BouncedCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=3", r: &g.Record.ComplaintCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=4", r: &g.Record.SuppressedCount},
	}

	for i := 0; i < len(gstat); i++ {
		q := fmt.Sprintf(gstat[i].q, d.ID)
		(*gstat[i].r), err = db.GetJoinSetCount(q)
		if err != nil {
			SvcGridErrorReturn(w, err)
			return
		}
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// GetGroup returns the requested group
// wsdoc {
//  @Title  Get Group
//	@URL /v1/getroup/GID
//  @Method  GET
//	@Synopsis Get information on a Group
//  @Description
//	@Input
//  @Response GroupGetResponse
// wsdoc }
func getGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getGroup"
	util.Console("entered %s\n", funcname)
	var g GroupGetResponse
	a, err := db.GetGroup(d.ID)
	if err != nil {
		SvcGridErrorReturn(w, err)
		return
	}
	if a.GID > 0 {
		var gg GroupGrid
		util.MigrateStructVals(&a, &gg)
		g.Record = gg
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
