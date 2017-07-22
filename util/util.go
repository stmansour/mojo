package util

import (
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

// DATEFMT et al., are commonly used date formats
const (
	DATEFMT        = "01/02/06"
	DATEFMT2       = "1/2/06"
	DATEFMT3       = "1/2/2006"
	DATEFMT4       = "01/02/2006"
	DATEINPFMT     = "2006-01-02"
	DATEFMTSQL     = DATEINPFMT
	DATETIMEINPFMT = "2006-01-02 15:04:00 MST"
	DATETIMEFMT    = "2006-01-02T15:04:00Z"
	DATEREPORTFMT  = "Jan 2, 2006"
)

// GenerateOptOutCode generates a reproducable code for the user. This code
// can be used to validate an opt-out link.
func GenerateOptOutCode(fn, ln, email string, pid int64) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s %d %s %s", fn, pid, email, ln))))
}

// ErrCheck - saves a bunch of typing, prints error if it exists
//            and provides a traceback as well
func ErrCheck(err error) {
	if err != nil {
		fmt.Printf("error = %v\n", err)
		debug.PrintStack()
		log.Fatal(err)
	}
}

// Stripchars removes the characters from chr in str and returns the updated string.
func Stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

// RemoveBackslash removes the backslash from in front of
// the double quote mark in byte buffers in returns a new
// byte buffer.  AWS SNS sends data this way and it breaks
// in the JSON decoder.  This routine fixes things up.
func RemoveBackslash(dat []byte) []byte {
	l := len(dat)
	b := make([]byte, l)
	if l == 0 {
		return b
	}
	lim := l - 1
	j := 0
	for i := 0; i < lim; i++ {
		if dat[i] == '\\' && dat[i+1] == '"' {
			continue // just skip the backslash
		}
		b[j] = dat[i]
		j++
	}
	b[j] = dat[lim]
	return b
}

// ScrubEmailAddr removes characters that are not allowed in an email address
// from the provided string and returns the updated string.
// It also removes some of the illegal character combinations I have encountered
// when doing large email blasts with Mojo
func ScrubEmailAddr(ss string) string {
	s1 := Stripchars(ss, " ,\"():;<>")

	// Fix addresses of the form: Dennis.E..Echelberry@faa.gov
	j := strings.Index(s1, "..")
	if j >= 0 {
		s := s1[:j+1]
		if len(s1) > j+2 {
			s += s1[j+2:]
		}
		s1 = s
	}

	// Fix addresses of the form: Richard.D.AndersonJr.@faa.gov
	if j = strings.Index(s1, ".@"); j >= 0 {
		s := s1[:j]
		if len(s1) > j+1 {
			s += s1[j+1:]
		}
		s1 = s
	}

	// Fix addresses of the form: I@mw@lkingw1th477@ng3ls.M.McDonald@faa.gov
	// save the rightmost @, just remove all other @ characters
	if strings.Count(s1, "@") > 1 {
		i := strings.LastIndex(s1, "@")
		sa := Stripchars(s1[:i], "@")
		sb := s1[i:]
		s1 = sa + sb
	}

	return s1
}

// LoadCSV loads a comma-separated-value file into an array of strings and returns the array of strings
func LoadCSV(fname string) [][]string {
	t := [][]string{}
	f, err := os.Open(fname)
	if nil == err {
		defer f.Close()
		reader := csv.NewReader(f)
		reader.FieldsPerRecord = -1
		reader.LazyQuotes = true
		rawCSVdata, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		t = append(t, rawCSVdata...)
	} else {
		Ulog("LoadCSV: could not open CSV file. err = %v\n", err)
	}
	return t
}

// Ulog is the standard logger
func Ulog(format string, a ...interface{}) {
	p := fmt.Sprintf(format, a...)
	log.Print(p)
	// debug.PrintStack()
}

// UlogAndPrint logs and prints the message
func UlogAndPrint(format string, a ...interface{}) {
	p := fmt.Sprintf(format, a...)
	log.Print(p)
	fmt.Print(p)
	// debug.PrintStack()
}

// LogAndPrintError encapsulates logging and printing an error.
// Note that the error is printed only if the environment is NOT production.
func LogAndPrintError(funcname string, err error) {
	errmsg := fmt.Sprintf("%s: err = %v\n", funcname, err)
	Ulog(errmsg)
	fmt.Println(errmsg)
}

// Tline returns a string of dashes that is the specified length
func Tline(n int) string {
	p := make([]byte, n)
	for i := 0; i < n; i++ {
		p[i] = '-'
	}
	return string(p)
}

// Mkstr returns a string of n of the supplied character that is the specified length
func Mkstr(n int, c byte) string {
	p := make([]byte, n)
	for i := 0; i < n; i++ {
		p[i] = c
	}
	return string(p)
}

// IsSQLNoResultsError returns true if the error provided is a sql err indicating no rows in the solution set.
func IsSQLNoResultsError(err error) bool {
	s := fmt.Sprintf("%v", err)
	return strings.Contains(s, "no rows in result")
}

// IntFromString converts the supplied string to an int64 value. If there
// is a problem in the conversion, it generates an error message. To suppress
// the error message, pass in "" for errmsg.
func IntFromString(sa string, errmsg string) (int64, error) {
	var n = int64(0)
	s := strings.TrimSpace(sa)
	if len(s) > 0 {
		i, err := strconv.Atoi(s)
		if err != nil {
			if "" != errmsg {
				return 0, fmt.Errorf("IntFromString: %s: %s", errmsg, s)
			}
			return n, err
		}
		n = int64(i)
	}
	return n, nil
}
