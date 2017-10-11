package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"mojo/db"
	"mojo/util"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var re = regexp.MustCompile(">([^<]+)<")
var reprof = regexp.MustCompile("a href=\"(\\(LoadPerson\\)[^\"]+)")
var urlbase = string("https://directory.faa.gov/appsPub/National/employeedirectory/faadir.nsf/")

// FAAScraper is a struct with context information for the FAA Scraper
var FAAScraper struct {
	GID     int64
	c       chan string
	quick   bool
	workers int
}

// InitFAAScraper is used to provide context information for this scraper
func InitFAAScraper(gid int64, q bool, w int) {
	FAAScraper.GID = gid
	FAAScraper.quick = q
	FAAScraper.workers = w

}

func safeMatrixGet(m [][]string, row, col int) string {
	if len(m) > row {
		if len(m[row]) > col {
			return strings.TrimSpace(m[row][col])
		}
	}
	return ""
}

// parseProfileMatrix takes a matrix of strings producted by parseProfileHTML and plucks
// out useful values.  The profile page consists of 42 lines of html. These are parsed into
// a matrix of 42 rows and up to 5 columns per row.  Some locations (row,column) in this
// matrix have useful information.  Here are the locations identified:
//
// 		Functional Job Title: (18,1)
// 		Service Unit: (19,1)
// 		Directoriate: (20,1)
//		Division: (21,1)
//		Office Phone: (23,1)
//		Physical Address:  (29-31,1)
//		Mail Address: (29-31,3)
// 		RoomNumber: (32,3)
// 		MailStop: (33,3)
//
func parseProfileMatrix(m [][]string, p *db.Person) {
	// Name: (16,0)   that is, row 16 column 0
	var first, middle, last string
	parseProfileName(m[16][0], &first, &middle, &last)
	// fmt.Printf("Name parsed as: first: %s, middle: %s, last: %s\n", first, middle, last)

	// if the first names do not match, assume that the smaller of the two is the Preferred Name
	if first != p.FirstName {
		fmt.Printf("Search FirstName does not match Profile Search name: %s vs %s\n", p.FirstName, first)
		if len(first) > len(p.FirstName) {
			p.PreferredName = p.FirstName
			p.FirstName = first
		} else {
			p.PreferredName = first
		}
	}
	if middle != p.MiddleName {
		fmt.Printf("Search MiddleName does not match Profile MiddleName: %s vs %s\n", p.MiddleName, middle)
	}
	if last != p.LastName {
		fmt.Printf("Search LastName does not match Profile LastName: %s vs %s\n", p.LastName, last)
	}

	phone := safeMatrixGet(m, 23, 1)
	if phone != p.OfficePhone {
		// look for innocuous miscompares, for example:  N/A  vs  ""
		if p.OfficePhone != "N/A" || len(phone) != 0 {
			// could be a problem
			fmt.Printf("Search OfficePhone does not match Profile OfficePhone: %s vs %s\n", p.OfficePhone, phone)
		}
	}

	// room number is at row 18. Mailstop is row 19
	p.RoomNumber = safeMatrixGet(m, 32, 3)
	p.MailStop = safeMatrixGet(m, 33, 3)

	// Update the address.  The address is contained in the array of strings between
	// rows 14-17 in column 3.   that is [14..17][3]
	var addr []string
	for i := 29; i <= 31; i++ { // these are the rows in which an address MIGHT appear
		if len(m[i]) >= 4 { // if there are 4 columns
			if len(m[i][3]) > 0 { // if there's anything in the column...
				addr = append(addr, m[i][3]) // grab it
			}
		}
	}
	parseCityStateZip(addr, p)

	if false { // debug printing
		fmt.Printf("Mail Address:    %s\n", p.MailAddress)
		fmt.Printf("Mail City:       %s\n", p.MailCity)
		fmt.Printf("Mail State:      %s\n", p.MailState)
		fmt.Printf("Mail PostalCode: %s\n", p.MailPostalCode)
	}
}

// removeTags returns a string in which all tags within the supplied string have been removed
// Playground: https://play.golang.org/p/p6Cvfxbn9j
func removeTags(s string) string {
	i := 0
	j := 0
	r := ""
	if i = strings.Index(s, "<"); i < 0 {
		return s
	}
	for {
		if j = strings.Index(s, ">"); j < 0 {
			return r
		}
		if j >= len(s)-1 {
			return r
		}
		s = s[j+1:]
		if s[0] == '<' {
			continue
		}
		if i = strings.Index(s, "<"); i < 0 {
			r += s
			return r
		}
		r += s[:i]
		s = s[i:]
	}
}

// extract the text out of each cell in a row and return a slice of strings
// Playground: https://play.golang.org/p/pvpktzDA98
func parseCellData(s string) []string {
	var sa []string
	j := -1 // index of cell within a row (<td> within a <tr>)
	i := 0  // initial index
	for i >= 0 {
		i = strings.Index(s, "<td")
		if i < 0 { // no more <td> cells in this line
			break
		}
		sa = append(sa, "") // we found *something*
		j++                 // index of what we just added to sa
		for i >= 0 {        // search for the text after the close of this <td>
			s = s[i+1:]                         // one character past '<'
			iFirstChar := strings.Index(s, ">") // close of the "<td" tag
			if iFirstChar < 0 {                 // if we don't find one, this is an unexpected line format
				break // badly formed line, break out and keep going
			}
			iFirstChar++                           // move past the '>'
			iLastChar := strings.Index(s, "</td>") // find the closing
			if iLastChar < 0 {                     // if we don't find one, this is an unexpected line format
				break // just break out of the loop and keep looking
			}
			if iLastChar > iFirstChar {
				sa[j] = removeTags(s[iFirstChar:iLastChar]) // grab the cell contents
			}
			s = s[iLastChar+5:] // new string to parse
			i = 0               // begin parsing here
			break
		}
	}
	return sa
}

// parseProfileHTML downloads the profile page for the person and extracts
// the name, address and room number.  The name may be different (Robert vs Bob).
// Typically the longer of the two names is the Formal name.  We'll use the longer
// of the two for creating the email address. The shorter one will become the
// PreferredName
//-----------------------------------------------------------------------------
func parseProfileHTML(s string, p *db.Person) {
	//fmt.Printf("Entered parseProfileHTML for: %s %s %s\n", p.FirstName, p.MiddleName, p.LastName)
	url := urlbase + s

	// let's do retries here... 3 tries.  Wait 5 seconds between each try...
	var err error
	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = http.Get(url)
		if nil == err {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		fmt.Printf("http.Get(%s) failed 3 times\n", url)
		fmt.Printf("\terr = %v\n", err)
		return // let's let the program keep running
	}

	defer resp.Body.Close()
	var body []byte
	for i := 0; i < 3; i++ {
		body, err = ioutil.ReadAll(resp.Body)
		if nil == err {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		fmt.Printf("loadProfile: ioutil.ReadAll failed 3 times\n")
		fmt.Printf("\terr = %v\n", err)
		return // let's let the program keep running
	}
	// fmt.Printf("Profile = %s\n", string(body))

	var strMatrix [][]string
	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	for i := 0; scanner.Scan(); i++ {
		s := scanner.Text()
		strMatrix = append(strMatrix, parseCellData(s))
	}
	parseProfileMatrix(strMatrix, p)

	if false { // for debugging
		for i := 0; i < len(strMatrix); i++ {
			fmt.Printf("%d.  ", i)
			for j := 0; j < len(strMatrix[i]); j++ {
				if len(strMatrix[i][j]) > 0 {
					fmt.Printf("[%d] %q    ", j, strMatrix[i][j])
				}
			}
			fmt.Println()
		}
	}
}

// MergePerson merges fields from A to B. It does NOT overwrite B's
// opt-out status, and it does NOT overwrite email addresses.
//-----------------------------------------------------------------------------
func MergePerson(a, b *db.Person) {
	b.FirstName = a.FirstName
	b.MiddleName = a.MiddleName
	b.LastName = a.LastName
	b.OfficePhone = a.OfficePhone
	b.MailCity = a.MailCity
	b.MailCountry = a.MailCountry
	b.MailState = a.MailState
	b.MailStop = a.MailStop
	b.JobTitle = a.JobTitle
	b.PreferredName = a.PreferredName
	b.RoomNumber = a.RoomNumber
}

// AddPersonToGroup creates a PGroup record for the specified pid,gid pair
// if it does not already exist.
func AddPersonToGroup(pid, gid int64) error {
	// see if they already exist...
	_, err := db.GetPGroup(pid, gid)
	if util.IsSQLNoResultsError(err) {
		var a = db.PGroup{PID: pid, GID: gid}
		// fmt.Printf("InsertPGroup:  pid = %d, gid = %d\n", pid, gid)
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

// InsertOrUpdatePerson looks in the database for a person that matches
// person p.  If found, fields are updated as appropriate and the record
// is updated.  If not found, person p will be inserted.
//-----------------------------------------------------------------------------
func InsertOrUpdatePerson(p *db.Person) {
	// Job title may have been updated.
	plist, err := db.GetPersonByRecordFieldMatching(p.FirstName, p.MiddleName, p.LastName, p.OfficePhone, p.MailAddress)
	if len(plist) > 1 {
		util.Ulog("ERROR: YIPES! %d people returned by GetPersonByRecordFieldMatching for name %s %s %s\n", len(plist), p.FirstName, p.MiddleName, p.LastName)
		return
	}
	if len(plist) == 1 {
		MergePerson(p, &plist[0])
		db.UpdatePerson(&plist[0]) // force the timestamp to update
		AddPersonToGroup(plist[0].PID, FAAScraper.GID)
		return
	}

	// see if there's a match on email...
	p2, err := db.GetPersonByEmail(p.Email1, p.Email2)
	if err != nil && !util.IsSQLNoResultsError(err) {
		util.Ulog("GetPersonByEmail returned: %s\n", err.Error())
		return
	}
	if err == nil && p2.PID > 0 && len(plist) > 0 {
		MergePerson(&p2, &plist[0])
		db.UpdatePerson(&plist[0]) // force the timestamp to update
		AddPersonToGroup(plist[0].PID, FAAScraper.GID)
		return
	}

	// no match on anything we trust, let's assume it's a new record
	err = db.InsertPerson(p)
	if err != nil {
		util.Ulog("db.InsertPerson returned: %s\n", err.Error())
		return
	}
	AddPersonToGroup(p.PID, FAAScraper.GID)
}

// parseHTML parses a line of the HTML file returned from the search tool on
// the FAA website. // Lines are in the following format (when the HTML tags
// are removed):
//
//    name, jobtitle, officephone, profile, orgchart
//    "Aakre, Dave C","ATSS","701-451-6805"," View Profile","N/A"
//
// profile and orgchart are just text in this file (http links removed), so
// ignore them. Go Playground: https://play.golang.org/p/-NGrjQ1A7p
//-----------------------------------------------------------------------------
func parseHTML(s string) {
	// fmt.Printf("entered parseHTML.  s = %s\n", s)
	r := re.FindAllStringSubmatch(s, -1)
	l := len(r)
	// fmt.Printf("submatches found %d\n", l)
	if l == 0 {
		return
	}
	var p db.Person
	if l >= 3 {
		name := strings.TrimSpace(r[0][1])
		if len(name) == 0 {
			return
		}
		nameHandler(name, &p)
		p.JobTitle = r[1][1]
		p.OfficePhone = r[2][1]
		profre := reprof.FindAllStringSubmatch(s, -1)
		if len(profre) > 0 {
			// fmt.Printf("profile url = %s\n", profre[0][1])
			parseProfileHTML(profre[0][1], &p)
		}
		emailBuilder(&p)
		InsertOrUpdatePerson(&p)
	}
}

func processSearchResults(q string) {
	URL := "https://directory.faa.gov/appsPub/National/employeedirectory/faadir.nsf/SearchForm?OpenForm"
	hc := http.Client{}

	form := url.Values{}
	form.Add("__Click", "862570240055C5F3.c191ad9beca4086705256f6b00650208/$Body/0.1158")
	form.Add("FAP_LastName", fmt.Sprintf("%s*", q))
	form.Add("FAP_FirstName", "")
	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(form.Encode()))
	util.ErrCheck(err)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Cookie", "BIGipServerpool_prd_directory.faa.gov_https=3200430747.47873.0000")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36")
	req.Header.Add("Referer", "https://directory.faa.gov/appsPub/National/employeedirectory/faadir.nsf/SearchForm?OpenForm")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("Origin", "https://directory.faa.gov")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Host", "directory.faa.gov")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	if false {
		dump, err := httputil.DumpRequest(req, false)
		util.ErrCheck(err)
		fmt.Printf("\n\ndumpRequest = %s\n", string(dump))
	}

	resp, err := hc.Do(req)
	util.ErrCheck(err)
	defer resp.Body.Close()

	if false {
		dump, err := httputil.DumpResponse(resp, true)
		util.ErrCheck(err)
		fmt.Printf("\n\ndumpResponse = %s\n", string(dump))
	}

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("error with ioutil.ReadAll:  %s\n", err.Error())
		return
	}
	bodyString := string(bodyBytes)

	scanner := bufio.NewScanner(strings.NewReader(bodyString))
	for scanner.Scan() {
		s := scanner.Text()
		parseHTML(s)
	}
}

// worker waits for a work item (string s) to come to it via the
// channel string. When it gets one, it calls processSearchResults to
// handle that string. It will continue doing this as long as more
// work is available via channel n.  Once n is closed, it will exit
// which invokes the deferred work group exit.
//---------------------------------------------------------------------
func worker(n chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for s := range n {
		processSearchResults(s)
	}
}

// ScrapeFAA scrapes the FAA directory website for information about
// its employees
//---------------------------------------------------------------------
func ScrapeFAA() {
	var du db.DataUpdate
	du.GID = FAAScraper.GID
	du.DtStart = time.Now()
	db.InsertDataUpdate(&du)

	//------------------------------------------
	// Create a pool of worker goroutines...
	//------------------------------------------
	FAAScraper.c = make(chan string)
	wg := new(sync.WaitGroup)
	fmt.Printf("Adding %d workers\n", FAAScraper.workers)
	for i := 0; i < FAAScraper.workers; i++ {
		wg.Add(1)
		go worker(FAAScraper.c, wg)
	}

	for i := 'a'; i <= 'z'; i++ {
		for j := 'a'; j <= 'z'; j++ {
			q := fmt.Sprintf("%c%c", i, j)
			fmt.Printf("Searching %c%c\n", i, j)
			FAAScraper.c <- q
			if FAAScraper.quick {
				break
			}
		}
		if FAAScraper.quick {
			break
		}
	}
	close(FAAScraper.c)

	// now just wait for the workers to finish everything...
	wg.Wait()
	du.DtStop = time.Now()
	elapsed := du.DtStop.Sub(du.DtStart)
	fmt.Printf("Elapsed time: %s\n", elapsed)
	err := db.UpdateDataUpdate(&du)
	if err != nil {
		util.Ulog("Error updating DataUpdate record: %s\n", err.Error())
	}
	g, err := db.GetGroupByName("FAA")
	if err != nil {
		util.Ulog("Error getting group FAA: %s\n", err.Error())
	}
	g.DtStop = du.DtStop
	err = db.UpdateGroup(&g)
	if err != nil {
		util.Ulog("Error inserting group: %s\n", err.Error())
	}
}
