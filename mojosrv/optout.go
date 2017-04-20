package main

import (
	"fmt"
	"io"
	"log"
	"mojo/db"
	"mojo/util"
	"net/http"
	"os"
	"time"

	"github.com/kardianos/osext"
)

// SendFileReply copies the supplied file to the output io.Writer w.
func SendFileReply(w io.Writer, fname string) {
	srcFile, err := os.Open(fname)
	if err != nil {
		util.Ulog("Error with os.Open: %s\n", err.Error())
	}
	defer srcFile.Close()
	_, err = io.Copy(w, srcFile) // check first var for number of bytes copied
	if err != nil {
		util.Ulog("Error with file io.Copy: %s\n", err.Error())
	}
}

// SvcOptOut returns the number of people in the database
// wsdoc {
//  @Title  Opt Out
//	@URL /v1/optout?e=emailaddr&c=code
//  @Method   GET
//	@Synopsis Opt-out of the mailings
//  @Descr  Sets the associated person Status to 1 indicating opt-out
//	@Input n/a
//  @Response page indicating the success or failure of the opt out
// wsdoc }
func SvcOptOut(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcOptOut"
	fmt.Printf("Entered %s\n", funcname)
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	q := r.URL.Query()
	email := q.Get("e")
	code := q.Get("c")

	fmt.Printf("Found email = %s, code = %s\n", email, code)

	p, err := db.GetPersonByEmail(email)
	if err != nil {
		fmt.Printf("EGetPersonByEmail %s returned:  %s", email, err.Error())
		SvcGridErrorReturn(w, err)
		return
	}

	s := util.GenerateOptOutCode(p.FirstName, p.LastName, p.Email1, p.PID)
	if s == code {
		fmt.Printf("Code confirmed, setting OptOut status\n")
		p.Status = db.OPTOUT
		p.OptOutDate = time.Now()
		err = db.UpdatePerson(&p)
		if err != nil {
			fmt.Printf("EGetPersonByEmail %s returned:  %s", email, err.Error())
			SvcGridErrorReturn(w, err)
			return
		}
		fmt.Printf("OptOut succeeded - return page\n")
		SendFileReply(w, folderPath+"/html/optouts.html")
		return
	}
	fmt.Printf("Code for %s should be %s, return error page\n", email, s)
	SendFileReply(w, folderPath+"/html/optoutf.html")

}
