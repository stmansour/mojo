package mailsend

import (
	"mojo/db"
	"mojo/util"
)

// AddPersonToGroup creates a PGroup record for the specified pid,gid pair
// if it does not already exist.
//-----------------------------------------------------------------------------
func AddPersonToGroup(pid, gid int64) error {
	// see if they already exist...
	_, err := db.GetPGroup(pid, gid)
	if util.IsSQLNoResultsError(err) {
		var a = db.PGroup{PID: pid, GID: gid}
		err = db.InsertPGroup(&a)
		if err != nil {
			util.Ulog("Error with InsertPGroup: %s\n", err.Error())
		}
		return err
	}
	if err == nil {
		return nil // they're already in the group
	}
	util.Ulog("Error trying to GetPGroup = %s\n", err.Error())
	return err
}

// RemovePersonFromGroup removes the suppled pid from the group
func RemovePersonFromGroup(pid, gid int64) error {
	err := db.DeletePGroup(pid, gid)
	if err != nil {
		if util.IsSQLNoResultsError(err) {
			return nil // this is fine
		}
		return err
	}
	return nil
}
