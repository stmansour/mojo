package main

// These are general utilty routines to support w2ui grid components.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mojo/util"
	"net/http"
	"net/url"
	"strings"
)

// SvcGridError is the generalized error structure to return errors to the grid widget
type SvcGridError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SvcStatusResponse is the response to return status when no other data
// needs to be returned
type SvcStatusResponse struct {
	Status string `json:"status"` // typically "success"
	Recid  int64  `json:"recid"`  // set to id of newly inserted record
}

// ServiceHandler describes the handler for all services
type ServiceHandler struct {
	Cmd     string
	Handler func(http.ResponseWriter, *http.Request, *ServiceData)
}

// GenSearch describes a search condition
type GenSearch struct {
	Field    string `json:"field"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

// ColSort is what the UI uses to indicate how the return values should be sorted
type ColSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// WebGridSearchRequestJSON is a struct suitable for describing a webservice operation.
// It is the wire format data. It will be merged into another object where JSONTime values
// are converted to time.Time
type WebGridSearchRequestJSON struct {
	Cmd         string      `json:"cmd"`         // get, save, delete
	Limit       int         `json:"limit"`       // max number to return
	Offset      int         `json:"offset"`      // solution set offset
	Selected    []int       `json:"selected"`    // selected rows
	SearchLogic string      `json:"searchLogic"` // OR | AND
	Search      []GenSearch `json:"search"`      // what fields and what values
	Sort        []ColSort   `json:"sort"`        // sort criteria
}

// WebGridSearchRequest is a struct suitable for describing a webservice operation.
type WebGridSearchRequest struct {
	Cmd         string      `json:"cmd"`         // get, save, delete
	Limit       int         `json:"limit"`       // max number to return
	Offset      int         `json:"offset"`      // solution set offset
	Selected    []int       `json:"selected"`    // selected rows
	SearchLogic string      `json:"searchLogic"` // OR | AND
	Search      []GenSearch `json:"search"`      // what fields and what values
	Sort        []ColSort   `json:"sort"`        // sort criteria
}

// WebFormRequest is a struct suitable for describing a webservice operation.
type WebFormRequest struct {
	Cmd      string      `json:"cmd"`    // get, save, delete
	Recid    int         `json:"recid"`  // max number to return
	FormName string      `json:"name"`   // solution set offset
	Record   interface{} `json:"record"` // selected rows
}

// WebTypeDownRequest is a search call made by a client while the user is
// typing in something to search for and the expecation is that the solution
// set will be sent back in realtime to aid the user.  Search is a string
// to search for -- it's what the user types in.  Max is the maximum number
// of matches to return.
type WebTypeDownRequest struct {
	Search string `json:"search"`
	Max    int    `json:"max"`
}

// WebGridDelete is a generic command structure returned when records are
// deleted from a grid. the Selected struct will contain the list of ids
// (recids which should map to the record type unique identifier) that are
// to be deleted.
type WebGridDelete struct {
	Cmd      string  `json:"cmd"`
	Selected []int64 `json:"selected"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
}

// ServiceData is the generalized data gatherer for svcHandler. It allows all
// the common data to be centrally parsed and passed to a handler, which may
// need to parse further to get its unique data.  It includes fields for
// common data elements in web svc requests
type ServiceData struct { // position 0 is 'v1'
	Service       string               // the service requested (position 1)
	GID           int64                // which group (position 2)
	ID            int64                // the numeric id parsed from position 3
	wsSearchReq   WebGridSearchRequest // what did the search requester ask for
	wsTypeDownReq WebTypeDownRequest   // fast for typedown
	data          string               // the raw unparsed data
	GetParams     map[string]string    // parameters when HTTP GET is used
}

// Svcs is the table of all service handlers
var Svcs = []ServiceHandler{
	{"group", SvcHandlerGroup},
	{"groups", SvcSearchHandlerGroups},
	{"groupcount", SvcGroupsCount},
	{"people", SvcSearchHandlerPeople},
	{"peoplecount", SvcPeopleCount},
	{"peoplestats", SvcPeopleStats},
	{"person", SvcHandlerPerson},
	{"ping", SvcHandlerPing},
	{"qrescount", SvcQueryResultsCount},
	{"queries", SvcSearchHandlerQueries},
	{"query", SvcHandlerQuery},
}

// SvcHandlerPing is the most basic test that you can run against the server
// see if it is alive and taking requests. It will return its version number.
func SvcHandlerPing(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	fmt.Fprintf(w, "Accord Mojo - Version %s\n", getVersionNo())
}

// V1ServiceHandler is the main dispatch point for WEB SERVICE requests
//
// The expected input is of the form:
//		request=%7B%22cmd%22%3A%22get%22%2C%22selected%22%3A%5B%5D%2C%22limit%22%3A100%2C%22offset%22%3A0%7D
// This is exactly what the w2ui grid sends as a request.
//
// Decoded, this message looks something like this:
//		request={"cmd":"get","selected":[],"limit":100,"offset":0}
//
// The leading "request=" is optional. This routine parses the basic information, then contacts an appropriate
// handler for more detailed processing.  It will set the Cmd member variable.
//
// W2UI sometimes sends requests that look like this: request=%7B%22search%22%3A%22s%22%2C%22max%22%3A250%7D
// using HTTP GET (rather than its more typical POST).  The command decodes to this: request={"search":"s","max":250}
//
//-----------------------------------------------------------------------------------------------------------
func V1ServiceHandler(w http.ResponseWriter, r *http.Request) {
	funcname := "V1ServiceHandler"
	svcDebugTxn(funcname, r)
	// var err error
	var d ServiceData

	d.ID = -1 // indicates it has not been set

	//-----------------------------------------------------------------------
	// pathElements:  0   1            2
	//               /v1/{subservice}/{ID}
	//-----------------------------------------------------------------------
	ss := strings.Split(r.RequestURI[1:], "?") // it could be GET command
	pathElements := strings.Split(ss[0], "/")
	lpe := len(pathElements)
	if lpe > 1 {
		d.Service = pathElements[1]
	}
	if lpe > 2 {
		var err error
		d.ID, err = util.IntFromString(pathElements[2], "bad ID")
		if err != nil {
			e := fmt.Errorf("ID in URL is invalid: %s", err.Error())
			fmt.Printf("***ERROR IN URL***  %s", e.Error())
			SvcGridErrorReturn(w, e)
			return
		}
	}

	svcDebugURL(r, &d)
	showRequestHeaders(r)

	switch r.Method {
	case "POST":
		if nil != getPOSTdata(w, r, &d) {
			return
		}
	case "GET":
		if nil != getGETdata(w, r, &d) {
			return
		}
	}

	showWebRequest(&d)

	//-----------------------------------------------------------------------
	//  Now call the appropriate handler to do the rest
	//-----------------------------------------------------------------------
	found := false
	for i := 0; i < len(Svcs); i++ {
		if Svcs[i].Cmd == d.Service {
			Svcs[i].Handler(w, r, &d)
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("**** YIPES! **** %s - Handler not found\n", r.RequestURI)
		e := fmt.Errorf("Service not recognized: %s", d.Service)
		fmt.Printf("***ERROR IN URL***  %s", e.Error())
		SvcGridErrorReturn(w, e)
		return
	}
	svcDebugTxnEnd()
}

// SvcGridErrorReturn formats an error return to the grid widget and sends it
func SvcGridErrorReturn(w http.ResponseWriter, err error) {
	var e SvcGridError
	e.Status = "error"
	e.Message = fmt.Sprintf("Error: %s\n", err.Error())
	b, _ := json.Marshal(e)
	SvcWrite(w, b)
}

// SvcGetInt64 tries to read an int64 value from the supplied string.
// If it fails for any reason, it sends writes an error message back
// to the caller and returns the error.  Otherwise, it returns an
// int64 and returns nil
func SvcGetInt64(s, errmsg string, w http.ResponseWriter) (int64, error) {
	i, err := util.IntFromString(s, "not an integer number")
	if err != nil {
		err = fmt.Errorf("%s: %s", errmsg, err.Error())
		SvcGridErrorReturn(w, err)
		return i, err
	}
	return i, nil
}

// SvcExtractIDFromURI extracts an int64 id value from position pos of the supplied uri.
// The URI is of the form returned by http.Request.RequestURI .  In particular:
//
//	pos:     0    1      2  3
//  uri:    /v1/rentable/34/421
//
// So, in the example uri above, a call where pos = 3 would return int64(421). errmsg
// is a string that will be used in the error message if the requested position had an
// error during conversion to int64. So in the example above, pos 3 is the RID, so
// errmsg would probably be set to "RID"
func SvcExtractIDFromURI(uri, errmsg string, pos int, w http.ResponseWriter) (int64, error) {
	var ID = int64(0)
	var err error

	sa := strings.Split(uri[1:], "/")
	// fmt.Printf("uri parts:  %v\n", sa)
	if len(sa) < pos+1 {
		err = fmt.Errorf("Expecting at least %d elements in URI: %s, but found only %d", pos+1, uri, len(sa))
		// fmt.Printf("err = %s\n", err)
		SvcGridErrorReturn(w, err)
		return ID, err
	}
	// fmt.Printf("sa[pos] = %s\n", sa[pos])
	ID, err = SvcGetInt64(sa[pos], errmsg, w)
	return ID, err
}

func getPOSTdata(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	funcname := "getPOSTdata"
	var err error
	htmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e := fmt.Errorf("%s: Error reading message Body: %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return e
	}
	fmt.Printf("\t- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("\thtmlData = %s\n", htmlData)
	if len(htmlData) == 0 {
		d.wsSearchReq.Cmd = "?"
		return nil
	}
	u, err := url.QueryUnescape(string(htmlData))
	if err != nil {
		e := fmt.Errorf("%s: Error with QueryUnescape: %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return e
	}
	fmt.Printf("\tUnescaped htmlData = %s\n", u)

	u = strings.TrimPrefix(u, "request=") // strip off "request=" if it is present
	d.data = u
	var wjs WebGridSearchRequestJSON
	err = json.Unmarshal([]byte(u), &wjs)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return e
	}
	util.MigrateStructVals(&wjs, &d.wsSearchReq)
	return err
}

func getGETdata(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	funcname := "getGETdata"
	s, err := url.QueryUnescape(strings.TrimSpace(r.URL.String()))
	if err != nil {
		e := fmt.Errorf("%s: Error with url.QueryUnescape:  %s", funcname, err.Error())
		SvcGridErrorReturn(w, e)
		return e
	}
	fmt.Printf("Unescaped query = %s\n", s)
	w2uiPrefix := "request="
	n := strings.Index(s, w2uiPrefix)
	fmt.Printf("n = %d\n", n)
	if n > 0 {
		fmt.Printf("Will process as Typedown\n")
		d.data = s[n+len(w2uiPrefix):]
		fmt.Printf("%s: will unmarshal: %s\n", funcname, d.data)
		if err = json.Unmarshal([]byte(d.data), &d.wsTypeDownReq); err != nil {
			e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
			SvcGridErrorReturn(w, e)
			return e
		}
		d.wsSearchReq.Cmd = "typedown"
	} else {
		fmt.Printf("Will process as web search command\n")
		d.wsSearchReq.Cmd = r.URL.Query().Get("cmd")
	}
	return nil
}

func showRequestHeaders(r *http.Request) {
	fmt.Printf("Headers:\n")
	for k, v := range r.Header {
		fmt.Printf("\t%s: ", k)
		for i := 0; i < len(v); i++ {
			fmt.Printf("%q  ", v[i])
		}
		fmt.Printf("\n")
	}
}

func showWebRequest(d *ServiceData) {
	if d.wsSearchReq.Cmd == "typedown" {
		fmt.Printf("Typedown:\n")
		fmt.Printf("\tSearch  = %q\n", d.wsTypeDownReq.Search)
		fmt.Printf("\tMax     = %d\n", d.wsTypeDownReq.Max)
	} else {
		fmt.Printf("\tSearchReq:\n")
		fmt.Printf("\t\tCmd           = %s\n", d.wsSearchReq.Cmd)
		fmt.Printf("\t\tLimit         = %d\n", d.wsSearchReq.Limit)
		fmt.Printf("\t\tOffset        = %d\n", d.wsSearchReq.Offset)
		fmt.Printf("\t\tsearchLogic   = %s\n", d.wsSearchReq.SearchLogic)
		for i := 0; i < len(d.wsSearchReq.Search); i++ {
			fmt.Printf("\t\tsearch[%d] - Field = %s,  Type = %s,  Value = %s,  Operator = %s\n", i, d.wsSearchReq.Search[i].Field, d.wsSearchReq.Search[i].Type, d.wsSearchReq.Search[i].Value, d.wsSearchReq.Search[i].Operator)
		}
		for i := 0; i < len(d.wsSearchReq.Sort); i++ {
			fmt.Printf("\t\tsort[%d] - Field = %s,  Direction = %s\n", i, d.wsSearchReq.Sort[i].Field, d.wsSearchReq.Sort[i].Direction)
		}
	}
}

func svcDebugTxn(funcname string, r *http.Request) {
	fmt.Printf("\n%s\n", util.Mkstr(80, '-'))
	fmt.Printf("URL:      %s\n", r.URL.String())
	fmt.Printf("METHOD:   %s\n", r.Method)
	fmt.Printf("Handler:  %s\n", funcname)
}

func svcDebugURL(r *http.Request, d *ServiceData) {
	//-----------------------------------------------------------------------
	// pathElements: 0         1     2
	// Break up {subservice}/{BUI}/{ID} into an array of strings
	// BID is common to nearly all commands
	//-----------------------------------------------------------------------
	ss := strings.Split(r.RequestURI[1:], "?") // it could be GET command
	pathElements := strings.Split(ss[0], "/")
	fmt.Printf("\t%s\n", r.URL.String()) // print before we strip it off
	for i := 0; i < len(pathElements); i++ {
		fmt.Printf("\t\t%d. %s\n", i, pathElements[i])
	}
}

func svcDebugTxnEnd() {
	fmt.Printf("END\n")
}

// SvcWriteResponse finishes the transaction with the W2UI client
func SvcWriteResponse(g interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(g)
	if err != nil {
		e := fmt.Errorf("Error marshaling json data: %s", err.Error())
		util.Ulog("SvcWriteResponse: %s\n", err.Error())
		SvcGridErrorReturn(w, e)
		return
	}
	SvcWrite(w, b)
}

// SvcWrite is a general write routine for service calls... it is a bottleneck
// where we can place debug statements as needed.
func SvcWrite(w http.ResponseWriter, b []byte) {
	fmt.Printf("first 200 chars of response: %-200.200s\n", string(b))
	// fmt.Printf("\nResponse Data:  %s\n\n", string(b))
	w.Write(b)
}

// SvcWriteSuccessResponse is used to complete a successful write operation on w2ui form save requests.
func SvcWriteSuccessResponse(w http.ResponseWriter) {
	var g = SvcStatusResponse{Status: "success"}
	w.Header().Set("Content-Type", "application/json")
	SvcWriteResponse(&g, w)
}

// SvcWriteSuccessResponseWithID is used to complete a successful write operation on w2ui form save requests.
func SvcWriteSuccessResponseWithID(w http.ResponseWriter, id int64) {
	var g = SvcStatusResponse{Status: "success", Recid: id}
	w.Header().Set("Content-Type", "application/json")
	SvcWriteResponse(&g, w)
}
