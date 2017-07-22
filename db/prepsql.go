package db

import (
	"database/sql"
	"mojo/util"
	"strings"
)

// PrepSQL is the collection of prepared sql statements
type PrepSQL struct {
	GetPerson                      *sql.Stmt
	GetPersonByEmail               *sql.Stmt
	GetAllPeopleOptedIn            *sql.Stmt
	InsertPerson                   *sql.Stmt
	UpdatePerson                   *sql.Stmt
	DeletePerson                   *sql.Stmt
	GetPersonByName                *sql.Stmt
	GetPersonByRecordFieldMatching *sql.Stmt
	GetGroup                       *sql.Stmt
	GetGroups                      *sql.Stmt
	GetGroupByName                 *sql.Stmt
	InsertGroup                    *sql.Stmt
	UpdateGroup                    *sql.Stmt
	DeleteGroup                    *sql.Stmt
	GetPGroup                      *sql.Stmt
	InsertPGroup                   *sql.Stmt
	DeletePGroup                   *sql.Stmt
	InsertDataUpdate               *sql.Stmt
	UpdateDataUpdate               *sql.Stmt
	GetDataUpdate                  *sql.Stmt
	GetDataUpdateByGroup           *sql.Stmt
	GetQuery                       *sql.Stmt
	GetQueryByName                 *sql.Stmt
	InsertQuery                    *sql.Stmt
	UpdateQuery                    *sql.Stmt
	DeleteQuery                    *sql.Stmt
}

// BuildPreparedStatements creates all the prepared statements for this db
func BuildPreparedStatements() {
	var err error
	var flds string
	var s1, s2, s3 string

	//--------------------------------------------
	//    PEOPLE
	//--------------------------------------------
	flds = "PID,FirstName,MiddleName,LastName,PreferredName,JobTitle,OfficePhone,OfficeFax,Email1,MailAddress,MailAddress2,MailCity,MailState,MailPostalCode,MailCountry,RoomNumber,MailStop,Status,OptOutDate,LastModTime,LastModBy"
	DB.DBFields["People"] = flds
	s1, s2, s3 = GenSQLInsertAndUpdateStrings(flds)
	DB.Prepstmt.GetPerson, err = DB.Db.Prepare("SELECT " + flds + " FROM People WHERE PID=?")
	util.ErrCheck(err)
	DB.Prepstmt.GetPersonByEmail, err = DB.Db.Prepare("SELECT " + flds + " FROM People WHERE Email1=?")
	util.ErrCheck(err)
	DB.Prepstmt.GetPersonByName, err = DB.Db.Prepare("SELECT " + flds + " FROM People WHERE FirstName=? AND MiddleName=? AND LastName=?")
	util.ErrCheck(err)
	DB.Prepstmt.GetPersonByRecordFieldMatching, err = DB.Db.Prepare("SELECT " + flds + " FROM People WHERE (FirstName=? OR PreferredName=?) AND MiddleName=? AND LastName=? and OfficePhone=? and MailAddress=?")
	util.ErrCheck(err)
	DB.Prepstmt.GetAllPeopleOptedIn, err = DB.Db.Prepare("SELECT People.* FROM People JOIN PGroup WHERE PGroup.GID=? AND Status=0 AND People.PID=PGroup.PID")
	util.ErrCheck(err)
	DB.Prepstmt.InsertPerson, err = DB.Db.Prepare("INSERT INTO People (" + s1 + ") VALUES(" + s2 + ")")
	util.ErrCheck(err)
	DB.Prepstmt.UpdatePerson, err = DB.Db.Prepare("UPDATE People SET " + s3 + " WHERE PID=?")
	util.ErrCheck(err)
	DB.Prepstmt.DeletePerson, err = DB.Db.Prepare("DELETE FROM People WHERE PID=?")
	util.ErrCheck(err)

	//--------------------------------------------
	//    EGROUPS - defines a group
	//--------------------------------------------
	flds = "GID,GroupName,GroupDescription,DtStart,DtStop,LastModTime,LastModBy"
	DB.DBFields["EGroup"] = flds
	s1, s2, s3 = GenSQLInsertAndUpdateStrings(flds)
	DB.Prepstmt.GetGroup, err = DB.Db.Prepare("SELECT " + flds + " FROM EGroup WHERE GID=?")
	util.ErrCheck(err)
	DB.Prepstmt.GetGroupByName, err = DB.Db.Prepare("SELECT " + flds + " FROM EGroup WHERE GroupName=?")
	util.ErrCheck(err)
	DB.Prepstmt.InsertGroup, err = DB.Db.Prepare("INSERT INTO EGroup (" + s1 + ") VALUES(" + s2 + ")")
	util.ErrCheck(err)
	DB.Prepstmt.UpdateGroup, err = DB.Db.Prepare("UPDATE EGroup SET " + s3 + " WHERE GID=?")
	util.ErrCheck(err)
	DB.Prepstmt.DeleteGroup, err = DB.Db.Prepare("DELETE FROM EGroup WHERE GID=?")
	util.ErrCheck(err)

	//--------------------------------------------------------------
	//    PGROUPS - defines the groups to which a person belongs
	//--------------------------------------------------------------
	flds = "PID,GID,LastModTime,LastModBy"
	DB.DBFields["PGroup"] = flds
	DB.Prepstmt.GetPGroup, err = DB.Db.Prepare("SELECT " + flds + " FROM PGroup WHERE PID=? AND GID=?")
	util.ErrCheck(err)
	DB.Prepstmt.InsertPGroup, err = DB.Db.Prepare("INSERT INTO PGroup (PID,GID,LastModBy) VALUES(?,?,?)")
	util.ErrCheck(err)
	DB.Prepstmt.DeletePGroup, err = DB.Db.Prepare("DELETE FROM PGroup WHERE PID=? AND GID=?")
	util.ErrCheck(err)

	//--------------------------------------------
	//    DATA UPDATE
	//--------------------------------------------
	flds = "DUID,GID,DtStart,DtStop,LastModTime,LastModBy"
	DB.DBFields["DataUpdate"] = flds
	s1, s2, s3 = GenSQLInsertAndUpdateStrings(flds)
	DB.Prepstmt.GetDataUpdate, err = DB.Db.Prepare("SELECT " + flds + " FROM DataUpdate WHERE DUID=?")
	util.ErrCheck(err)
	DB.Prepstmt.GetDataUpdateByGroup, err = DB.Db.Prepare("SELECT " + flds + " FROM DataUpdate WHERE GID=? ORDER BY DtStop DESC")
	util.ErrCheck(err)
	DB.Prepstmt.InsertDataUpdate, err = DB.Db.Prepare("INSERT INTO DataUpdate (" + s1 + ") VALUES(" + s2 + ")")
	util.ErrCheck(err)
	DB.Prepstmt.UpdateDataUpdate, err = DB.Db.Prepare("UPDATE DataUpdate SET " + s3 + " WHERE DUID=?")
	util.ErrCheck(err)

	//--------------------------------------------
	//    QUERY
	//--------------------------------------------
	flds = "QID,QueryName,QueryDescr,QueryJSON,LastModTime,LastModBy"
	DB.DBFields["Query"] = flds
	s1, s2, s3 = GenSQLInsertAndUpdateStrings(flds)
	DB.Prepstmt.GetQuery, err = DB.Db.Prepare("SELECT " + flds + " FROM Query WHERE QID=?")
	util.ErrCheck(err)
	DB.Prepstmt.GetQueryByName, err = DB.Db.Prepare("SELECT " + flds + " FROM Query WHERE QueryName=?")
	util.ErrCheck(err)
	DB.Prepstmt.InsertQuery, err = DB.Db.Prepare("INSERT INTO Query (" + s1 + ") VALUES(" + s2 + ")")
	util.ErrCheck(err)
	DB.Prepstmt.UpdateQuery, err = DB.Db.Prepare("UPDATE Query SET " + s3 + " WHERE QID=?")
	util.ErrCheck(err)
	DB.Prepstmt.DeleteQuery, err = DB.Db.Prepare("DELETE FROM Query WHERE QID=?")
	util.ErrCheck(err)
}

var mySQLRpl = "?"
var myRpl = mySQLRpl

// GenSQLInsertAndUpdateStrings generates a string suitable for SQL INSERT and UPDATE statements given the fields as used in SELECT statements.
//
//  example:
//	given this string:      "LID,BID,RAID,GLNumber,Status,Type,Name,AcctType,RAAssociated,LastModTime,LastModBy"
//  we return these five strings:
//  1)  "BID,RAID,GLNumber,Status,Type,Name,AcctType,RAAssociated,LastModBy"                  	-- use for SELECT
//  2)  "?,?,?,?,?,?,?,?,?"  																	-- use for INSERT
//  3)  "BID=?RAID=?,GLNumber=?,Status=?,Type=?,Name=?,AcctType=?,RAAssociated=?,LastModBy=?"   -- use for UPDATE
//
// Note that in this convention, we remove LastModTime from insert and update statements (the db is set up to update them by default) and
// we remove the initial ID as that number is AUTOINCREMENT on INSERTs and is not updated on UPDATE.
func GenSQLInsertAndUpdateStrings(s string) (string, string, string) {
	sa := strings.Split(s, ",")
	s2 := sa[1:]  // skip the ID
	l2 := len(s2) // how many fields
	if l2 > 2 {
		if s2[l2-2] == "LastModTime" { // if the last 2 values are "LastModTime" and "LastModBy"...
			s2[l2-2] = s2[l2-1] // ...move "LastModBy" to the previous slot...
			s2 = s2[:l2-1]      // ...and remove value .  We don't write LastModTime because it is set to automatically update
		}
	}
	s = strings.Join(s2, ",")
	l2 = len(s2) // may have changed

	// now s2 has the proper number of fields.  Produce a
	s3 := myRpl + ","               // start of the INSERT string  -- FOR USE WITH PRIMARY KEY AUTOINCREMENT
	s4 := s2[0] + "=" + myRpl + "," // start of the UPDATE string
	for i := 1; i < l2; i++ {
		s3 += myRpl               // for the INSERT string
		s4 += s2[i] + "=" + myRpl // for the UPDATE string
		if i < l2-1 {             // if there are more fields to come...
			s3 += "," // ...add a comma...
			s4 += "," // ...to both strings
		}
	}
	return s, s3, s4
}
