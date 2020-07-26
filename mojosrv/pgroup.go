package main

import (
	"fmt"
	"mojo/db"
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

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SavePGroup is sent to save one of open time slots as a reservation
type SavePGroup struct {
	Cmd    string     `json:"cmd"`
	Record PGroupItem `json:"record"`
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

func getPGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var g PGroupList
	q := fmt.Sprintf(`SELECT EGroup.GID,EGroup.GroupName FROM PGroup
INNER JOIN People ON (People.PID=PGroup.PID AND People.PID=%d)
INNER JOIN EGroup ON (EGroup.GID = PGroup.GID)
ORDER BY EGroup.GroupName ASC LIMIT %d OFFSET %d;`, d.ID, d.wsSearchReq.Limit, d.wsSearchReq.Offset)

	util.Console("query = %s\n", q)

	rows, err := db.DB.Db.Query(q)
	if err != nil {
		util.Console("Error from DB Query: %s\n", err.Error())
		SvcGridErrorReturn(w, err)
		return
	}
	defer rows.Close()

	i := int64(d.wsSearchReq.Offset)

	g.Total = 0
	for rows.Next() {
		var a PGroupItem
		if err = rows.Scan(&a.GID, &a.GroupName); err != nil {
			SvcGridErrorReturn(w, err)
			return
		}
		a.Recid = a.GID
		g.Records = append(g.Records, a)
		g.Total++ // update the g.Total  only after adding the record
		if g.Total >= int64(d.wsSearchReq.Limit) {
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

func savePGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {

}

func deletePGroup(w http.ResponseWriter, r *http.Request, d *ServiceData) {

}
