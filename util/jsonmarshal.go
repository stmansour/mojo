package util

import (
	"fmt"
	"strings"
	"time"
)

// RRDATEFMT is a shorthand date format used for text output
// Use these values:	Mon Jan 2 15:04:05 MST 2006
// const RRDATEFMT = "02-Jan-2006 3:04PM MST"
// const RRDATEFMT = "01/02/06 3:04PM MST"
const (
	RRDATEFMT        = "01/02/06"
	RRDATEFMT2       = "1/2/06"
	RRDATEFMT3       = "1/2/2006"
	RRDATEFMT4       = "01/02/2006"
	RRDATEINPFMT     = "2006-01-02"
	RRDATEFMTSQL     = RRDATEINPFMT
	RRDATETIMEINPFMT = "2006-01-02 15:04:00 MST"
	RRDATETIMEFMT    = "2006-01-02T15:04:00Z"
	RRDATEREPORTFMT  = "Jan 2, 2006"
)

var acceptedDateFmts = []string{
	RRDATEINPFMT,
	RRDATEFMT2,
	RRDATEFMT,
	RRDATEFMT3,
	RRDATETIMEINPFMT,
	RRDATETIMEFMT,
}

// StringToDate tries to convert the supplied string to a time.Time value. It will use the
// formats called out in dbtypes.go:  RRDATEFMT, RRDATEINPFMT, RRDATEINPFMT2, ...
//-----------------------------------------------------------------------------
func StringToDate(s string) (time.Time, error) {
	// try the ansi std date format first
	var Dt time.Time
	var err error
	s = strings.TrimSpace(s)
	for i := 0; i < len(acceptedDateFmts); i++ {
		Dt, err = time.Parse(acceptedDateFmts[i], s)
		if nil == err {
			break
		}
	}
	return Dt, err
}

// JSONDate is a wrapper around time.Time. We need it
// in order to be able to control the formatting used
// on the date values sent to the w2ui controls.  Without
// this wrapper, the default time format used by the
// JSON encoder / decoder does not work with the w2ui
// controls
type JSONDate time.Time

var earliestDate = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)

// MarshalJSON overrides the default time.Time handler and sends
// date strings of the form YYYY-MM-DD. Any date prior to Jan 1, 1900
// is snapped to Jan 1, 1900.
//-----------------------------------------------------------------------------
func (t *JSONDate) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t)
	if ts.Before(earliestDate) {
		ts = earliestDate
	}
	// val := fmt.Sprintf("\"%s\"", ts.Format("2006-01-02"))
	val := fmt.Sprintf("\"%s\"", ts.Format(RRDATEFMT3))
	return []byte(val), nil
}

// UnmarshalJSON overrides the default time.Time handler and reads in
// date strings of the form YYYY-MM-DD.  Any date prior to Jan 1, 1900
// is snapped to Jan 1, 1900.
//--------------------------------------------------------------------
func (t *JSONDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = Stripchars(s, "\"")
	// x, err := time.Parse("2006-01-02", s)
	x, err := StringToDate(s)
	if err != nil {
		return err
	}
	if x.Before(earliestDate) {
		x = earliestDate
	}
	*t = JSONDate(x)
	return nil
}

// JSONDateTime is a wrapper around time.Time. We need it
// in order to be able to control the formatting used
// on the DATETIME values sent to the w2ui controls.  Without
// this wrapper, the default time format used by the
// JSON encoder / decoder does not work with the w2ui
// controls
type JSONDateTime time.Time

// MarshalJSON overrides the default time.Time handler and sends
// date strings of the form YYYY-MM-DD. Any date prior to Jan 1, 1900
// is snapped to Jan 1, 1900.
//--------------------------------------------------------------------
func (t *JSONDateTime) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t)
	if ts.Before(earliestDate) {
		ts = earliestDate
	}
	val := fmt.Sprintf("\"%s\"", ts.Format(RRDATETIMEFMT))
	return []byte(val), nil
}

// UnmarshalJSON overrides the default time.Time handler and reads in
// date strings of the form YYYY-MM-DD.  Any date prior to Jan 1, 1900
// is snapped to Jan 1, 1900.
//--------------------------------------------------------------------
func (t *JSONDateTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = Stripchars(s, "\"")
	// x, err := time.Parse("2006-01-02", s)
	x, err := StringToDate(s)
	if err != nil {
		return err
	}
	if x.Before(earliestDate) {
		x = earliestDate
	}
	*t = JSONDateTime(x)
	return nil
}

// UnmarshalJSON a string to a W2uiHTMLSelect struct.  The data can come in
// two different forms.
//
//  1. 	{"id":"OKC","text":"OKC"}
// 		This is the expected form.  We just call unmarshal it.
//  2.  ""
//      When this happens it means that the dropdown was uninitialized and
//      the user didn't select anything. In this case just create a struct
//      and set the values to ""
//--------------------------------------------------------------------
func (t *W2uiHTMLSelect) UnmarshalJSON(b []byte) error {
	m := W2uiHTMLSelect{ID: "", Text: ""}
	// fmt.Printf("b: len = %d, contents: %s\n", len(b), string(b))
	if len(b) > 2 { // if b contains more than ""
		s := Stripchars(string(b), `"{}`) // {"id":"Person","text":"Person"}  -->  id:Person,text:Person
		ss := strings.Split(s, ",")       // ["id:Person" "text:Person"]
		sss := strings.Split(ss[0], ":")  // ["id" "Person"]
		m.ID = sss[1]                     // "Person"
		sss = strings.Split(ss[1], ":")   // ["text" "Person"]
		m.Text = sss[1]                   // "Person"
	}
	*t = m
	return nil
}

// MarshalJSON deals with XJSONYesNo
//--------------------------------------------------------------------
func (t *XJSONYesNo) MarshalJSON() ([]byte, error) {
	return []byte("\"" + string(*t) + "\""), nil
}

// UnmarshalJSON deals with XJSONYesNo
//--------------------------------------------------------------------
func (t *XJSONYesNo) UnmarshalJSON(b []byte) error {
	*t = XJSONYesNo(Stripchars(string(b), "\""))
	return nil
}
