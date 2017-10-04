package mailsend

import (
	"fmt"
	"mojo/db"
	"mojo/util"
	"rentroll/rlib"
)

// ValidateGroupEmailAddresses looks up the email address for every person
// in the supplied group and validates it.  If the email address validation
// fails, it prints out a message identifying the problem email address.
//
// INPUTS
// grp - name of the group to examine
//
// RETURNS
// error - any error encountered
//-----------------------------------------------------------------------------
func ValidateGroupEmailAddresses(grp string) error {
	var PID int64
	var Email1 string

	g, err := db.GetGroupByName(grp)
	if err != nil {
		return err
	}
	q := fmt.Sprintf("SELECT People.PID,People.Email1 FROM People INNER JOIN PGroup ON (PGroup.PID=People.PID AND PGroup.GID=%d)", g.GID)
	rows, err := db.DB.Db.Query(q)
	rlib.Errcheck(err)
	defer rows.Close()
	good := 0
	bad := 0
	empty := 0
	for rows.Next() {
		err := rows.Scan(&PID, &Email1)
		if err != nil {
			util.Ulog("Error with ReadPersonFromRows: %s\n", err.Error())
			return err
		}
		// fmt.Printf("Sending to %s\n", p.Email1)
		if len(Email1) == 0 {
			empty++
			continue
		}
		if !util.ValidEmailAddress(Email1) {
			bad++
			fmt.Printf("Invalid email address:  PID = %d, email = %q\n", PID, Email1)
		} else {
			good++
		}
	}
	fmt.Printf("Processed %d records\n", bad+good+empty)
	fmt.Printf("      Valid email addresses: %d\n", good)
	fmt.Printf("    Invalid email addresses: %d\n", bad)
	fmt.Printf("      Empty email addresses: %d\n", empty)

	return rows.Err()
}
