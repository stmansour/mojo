package main

import (
	"fmt"
	"mojo/db"
	"mojo/util"
	"rentroll/rlib"
	"strings"
)

// emailBuilder generates an email address based on the apparent
// default formula that the FAA uses for their email addresses.
// That is:
//		[firstName].[lastName]@FAA.gov
// or
//		[firstName].[middleInitial].[lastName]@FAA.gov
func emailBuilder(p *db.Person) {
	if len(p.MiddleName) > 0 {
		p.Email1 = util.ScrubEmailAddr(fmt.Sprintf("%s.%s.%s@faa.gov", p.FirstName, p.MiddleName, p.LastName))
	} else if len(p.FirstName) > 0 {
		p.Email1 = util.ScrubEmailAddr(fmt.Sprintf("%s.%s@faa.gov", p.FirstName, p.LastName))
	}
}

// parseProfileName - break up a name into first, middle, last
func parseProfileName(name string, f, m, l *string) {
	if len(name) > 0 {
		na := strings.Split(name, " ")
		switch len(na) {
		case 2:
			*f = na[0]
			*l = na[1]
		case 3:
			*f = na[0]
			*m = na[1]
			*l = na[2]
		default:
			fmt.Printf("unrecognized name format: %#v\n", na)
			return
		}
	}
}

func nameHandler(s string, p *db.Person) {
	// first, split last and first
	sa := strings.Split(s, ",")
	l := len(sa)
	for i := 0; i < l; i++ {
		sa[i] = strings.TrimSpace(sa[i])
	}

	// see if there is anything extra in the first name that we can
	// use as a middle name or initial
	if l == 2 {
		ta := strings.Split(sa[1], " ")
		if len(ta) > 1 {
			sa[1] = ta[0]
			sa = append(sa, ta[1])
			l = len(sa)
		}
	}
	switch {
	case l == 3:
		p.MiddleName = strings.TrimSpace(sa[2])
		fallthrough
	case l == 2:
		p.LastName = strings.TrimSpace(sa[0])
		p.FirstName = strings.TrimSpace(sa[1])
	case l == 1:
		p.LastName = strings.TrimSpace(sa[0])
	default:
		fmt.Printf("unknown format: sa = %#v\n", sa)
	}
}

// badAddress  if an addressed failed to parse, this routine is called to
// notify of the error.
//-----------------------------------------------------------------------------
func badAddress(s string, p *db.Person) {
	fmt.Printf("db.Person = %#v\n", p)
	if len(s) == 0 {
		rlib.Errcheck(fmt.Errorf("Address string s:  len(s) == 0"))
	}
	rlib.Errcheck(fmt.Errorf("Unrecognized address format:  %s", s))
}

// parseCityStateZip  does its best to pull a city, state, and zipcode out of
// the array of address strings we scraped from the profile page.
//-----------------------------------------------------------------------------
func parseCityStateZip(address []string, p *db.Person) {
	if len(address) == 0 {
		return
	}
	var s string
	s = strings.TrimSpace(address[len(address)-1])
	if len(s) == 0 {
		return
	}
	sa := strings.Split(s, ",")
	l := len(sa)
	if l == 2 {
		p.MailCity = strings.TrimSpace(sa[0])
		ta := strings.Split(sa[1], " ")
		if len(ta) > 1 {
			p.MailPostalCode = strings.TrimSpace(ta[len(ta)-1])
			p.MailState = strings.TrimSpace(strings.Join(ta[0:len(ta)-1], " "))
		}
	} else {
		// LOOK FOR KNOWN ERRONEOUS PATTERNS
		// Fort Worth, TX 76177, TX 76177
		if l == 3 {
			if strings.TrimSpace(sa[1]) == strings.TrimSpace(sa[2]) {
				p.MailCity = strings.TrimSpace(sa[0])
				ta := strings.Split(sa[1], " ")
				if len(ta) > 1 {
					p.MailPostalCode = strings.TrimSpace(ta[len(ta)-1])
					p.MailState = strings.TrimSpace(strings.Join(ta[0:len(ta)-1], " "))
				}
			} else {
				badAddress(s, p)
			}
		} else {
			badAddress(s, p)
		}
	}
	p.MailAddress = address[len(address)-2]
	return
}
