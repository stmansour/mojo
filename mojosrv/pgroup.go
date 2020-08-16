package main

import (
	"encoding/json"
	"fmt"
	"mojo/db"
	"mojo/mailsend"
	"mojo/util"
	"net/http"
)

// PGroup is the service for managing a persons group memeberships

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// PGroupItem describes an individual group to which the person belongs
type PGroupItem struct {
	Recid     int64 `json:"recid"`
	GID       int64
	GroupName string
}

// PGroupList is the full list of groups to which a person belongs
type PGroupList struct {
	Status  string       `json:"status"`
	Total   int64        `json:"total"`
	Records []PGroupItem `json:"records"`
}

// GroupTypeDown is the struct needed to match names in typedown controls
type GroupTypeDown struct {
	Recid     int64 `json:"recid"` // this will hold the RID
	GID       int64
	GroupName string
}

// GroupTypedownResponse is the data structure for the response to a search for people
type GroupTypedownResponse struct {
	Status  string          `json:"status"`
	Total   int64           `json:"total"`
	Records []GroupTypeDown `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SavePGroup is sent to save one of open time slots as a reservation
type SavePGroup struct {
	Cmd    string     `json:"cmd"`
	Record PGroupItem `json:"record"`
}

// GroupMembership holds an array with all the groups that a person
// currently belongs to.
type GroupMembership struct {
	Cmd    string `json:"cmd"`
	Groups []int64
}

//-----------------------------------------------------------------------------
//##########################################################################################################################################################
//-----------------------------------------------------------------------------

// SvcHandlerPGroup formats a complete data record for an assessment for use with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
// The server command can be:
//      get
//      save
//      delete
//-----------------------------------------------------------------------------------
func SvcHandlerPGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerPGroup\n")

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID < 0 {
			SvcGridErrorReturn(w, fmt.Errorf("PersonID is required but was not specified"))
			return
		}
		getPGroup(w, r, d)
		break
	case "save":
		savePGroup(w, r, d)
		break
	case "delete":
		deletePGroup(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcGridErrorReturn(w, err)
		return
	}
}

// SvcHandlerGroupMembership updates the groups that are associated with the PID
// cmd should be set to "save"
//
// a list of GIDs is passed in. This is the list of groups for the person.
// Some groups may need to be added, some deleted.
//
// URL:  /v1/groupmembership/PID
//       cmd: "save"
//       groups: [ 4, 6, 7]
//-----------------------------------------------------------------------------------
func SvcHandlerGroupMembership(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcHandlerGroupMembership"
	var err error
	util.Console("Entered %s\n", funcname)
	var a GroupMembership
	if err = json.Unmarshal([]byte(d.data), &a); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return
	}
	util.Console("Read %d group ids\n", len(a.Groups))

	// to what groups does this person currently belong?
	var gcur PGroupList
	if gcur, err = getPGroupList(w, r, d, d.ID, false); err != nil {
		SvcGridErrorReturn(w, err)
	}
	util.Console("Current group list for PID %d is:\n", d.ID)
	for i := 0; i < len(gcur.Records); i++ {
		util.Console("%s (%d)\n", gcur.Records[i].GroupName, gcur.Records[i].GID)
	}

	//------------------------------------------------------------------------
	// which groups do we need to add?  If we don't find a.Groups[i].GID in
	// the current list, then we need to add it.
	//------------------------------------------------------------------------
	for i := 0; i < len(a.Groups); i++ {
		found := false
		for j := 0; j < len(gcur.Records); j++ {
			if gcur.Records[j].GID == a.Groups[i] {
				found = true
				break
			}
		}
		if !found {
			// util.Console("Add to GID: %d\n", a.Groups[i])
			err = mailsend.AddPersonToGroup(d.ID, a.Groups[i])
		}
	}

	//------------------------------------------------------------------------
	// Remove the person from any group that does not appear in list...
	//------------------------------------------------------------------------
	for i := 0; i < len(gcur.Records); i++ {
		found := false
		for j := 0; j < len(a.Groups); j++ {
			if gcur.Records[i].GID == a.Groups[j] {
				found = true
				break
			}
		}
		if !found {
			util.Console("Remove from GID: %d\n", gcur.Records[i].GID)
			err = mailsend.RemovePersonFromGroup(d.ID, gcur.Records[i].GID)
		}
	}
	SvcWriteSuccessResponse(w)

}

// SvcGroupTD handles typedown messages when a person is looking for a group
// name
//-----------------------------------------------------------------------------
func SvcGroupTD(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "SvcGroupTD"
	var (
		g   GroupTypedownResponse
		m   []db.EGroup
		err error
	)
	util.Console("Entered %s\n", funcname)
	util.Console("handle typedown: GetGroupTypedown( search=%s, limit=%d\n", d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	m, err = db.GetGroupTypedown(r.Context(), d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	if err != nil {
		e := fmt.Errorf("Error getting typedown matches: %s", err.Error())
		SvcErrorReturn(w, e)
		return
	}

	for i := 0; i < len(m); i++ {
		var t GroupTypeDown
		t.GID = m[i].GID
		t.Recid = t.GID
		t.GroupName = m[i].GroupName
		g.Records = append(g.Records, t)
	}

	util.Console("GetRentableTypedown returned %d matches\n", len(g.Records))
	g.Total = int64(len(g.Records))
	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// getPGroupList returns a list of PGroupItems that is the list of groups to
// which the user currently belongs.
//
// INPUTS
//  w
//  t
//  d
//  PID -  the user of interest. May not be the same as d.ID
//  lim  - if you want the full list, no matter how many is in it, set this
//         to false.  Otherwise, when it's true, it will limit to the amount
//         specified in d.wsSearchReq.Limit and start at d.wsSearchReq.Offset
//-----------------------------------------------------------------------------
func getPGroupList(w http.ResponseWriter, r *http.Request, d *ServiceData, PID int64, lim bool) (PGroupList, error) {
	var g PGroupList

	q := fmt.Sprintf(`SELECT EGroup.GID,EGroup.GroupName FROM PGroup
INNER JOIN People ON (People.PID=PGroup.PID AND People.PID=%d)
INNER JOIN EGroup ON (EGroup.GID = PGroup.GID)
ORDER BY EGroup.GroupName ASC`, PID)

	if lim {
		q += fmt.Sprintf(` LIMIT %d OFFSET %d`, d.wsSearchReq.Limit, d.wsSearchReq.Offset)
	}

	q += ";"

	util.Console("query = %s\n", q)

	rows, err := db.DB.Db.Query(q)
	if err != nil {
		return g, err
	}
	defer rows.Close()

	i := int64(d.wsSearchReq.Offset)

	g.Total = 0
	for rows.Next() {
		var a PGroupItem
		if err = rows.Scan(&a.GID, &a.GroupName); err != nil {
			return g, err
		}
		a.Recid = a.GID
		g.Records = append(g.Records, a)
		g.Total++ // update the g.Total  only after adding the record
		if int64(d.wsSearchReq.Limit) > 0 && g.Total >= int64(d.wsSearchReq.Limit) {
			break // if we've added the max number requested, then exit
		}
		i++
	}
	util.Console("g.Total = %d\n", g.Total)
	util.ErrCheck(rows.Err())
	g.Status = "success"
	return g, nil
}

func getPGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var g PGroupList
	var err error
	if g, err = getPGroupList(w, r, d, d.ID, true); err != nil {
		SvcGridErrorReturn(w, err)
		return
	}
	SvcWriteResponse(&g, w)
}

func savePGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {

}

func deletePGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {

}
