package mailsend

import (
	"fmt"
	"mojo/db"
)

// GroupStats collects totals about the statuses of group members
type GroupStats struct {
	MemberCount     int64
	MailToCount     int64
	OptOutCount     int64
	BouncedCount    int64
	ComplaintCount  int64
	SuppressedCount int64
}

// GetGroupStats collects the status totals for members of the supplied group id
//
// INPUTS
//  id = GID of group
//
// RETURNS
//  GroupStats struct
//  any errors encountered
//-----------------------------------------------------------------------------
func GetGroupStats(id int64) (GroupStats, error) {
	var g GroupStats
	var err error
	var gstat = []struct {
		q string
		r *int64
	}{
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d", r: &g.MemberCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=0", r: &g.MailToCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=1", r: &g.OptOutCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=2", r: &g.BouncedCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=3", r: &g.ComplaintCount},
		{q: "select count(People.PID) FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=%d WHERE People.Status=4", r: &g.SuppressedCount},
	}

	for i := 0; i < len(gstat); i++ {
		q := fmt.Sprintf(gstat[i].q, id)
		(*gstat[i].r), err = db.GetJoinSetCount(q)
		if err != nil {
			return g, err
		}
	}
	return g, nil
}
