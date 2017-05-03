package main

import (
	"fmt"
	"mojo/db"
	"mojo/util"
	"rentroll/rlib"
	"strings"
)

// fixDoubleDotEmail fixes one error situation I found during the first email
// blast to the FAA.  Some people put a period after their middle initial in their name.
// This should be removed by the scraper.  But there were about 50 addresses that have
// the double dot already in the database.  This routine fixes those issues.
func fixDoubleDotEmail() {
	util.UlogAndPrint("Fix double-dot email addresses\n")
	q := "SELECT " + db.DB.DBFields["People"] + " FROM People WHERE Email1 LIKE \"%..%\""
	rows, err := db.DB.Db.Query(q)
	rlib.Errcheck(err)
	defer rows.Close()
	i := 0
	for rows.Next() {
		p, err := db.ReadPersonFromRows(rows)
		if err != nil {
			util.Ulog("Error with ReadPersonFromRows: %s\n", err.Error())
			return
		}
		fmt.Printf("Found: PID=%d, email = %s  (%s %s %s)\n", p.PID, p.Email1, p.FirstName, p.MiddleName, p.LastName)
		j := strings.Index(p.Email1, "..")
		s := p.Email1[:j+1]
		if len(p.Email1) > j+2 {
			s += p.Email1[j+2:]
		}
		fmt.Printf("Fixed address = %s\n", s)
		p.Email1 = s
		if err = db.UpdatePerson(&p); err != nil {
			util.Ulog("Error with UpdatePerson on PID=%d: %s\n", p.PID, err.Error())
			return
		}
		i++
	}
	fmt.Printf("Found %d double-dot email addresses\n", i)
	if i > 0 {
		fmt.Printf("Fixed all of them\n")
	}
}
