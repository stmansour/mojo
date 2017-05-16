package util

import "testing"

type emaildata struct {
	bad  string
	good string
}

// TestEmailScrub tests the conversion between int and XJSONYesNo
func TestEmailScrub(t *testing.T) {
	var m = []emaildata{
		{"sally@blob.com", "sally@blob.com"},
		{"Dennis.E..Echelberry@faa.gov", "Dennis.E.Echelberry@faa.gov"},
		{"Richard.D.AndersonJr.@faa.gov", "Richard.D.AndersonJr@faa.gov"},
		{"Billy.(Bob).Thorton@xyz.com", "Billy.Bob.Thorton@xyz.com"},
		{"I@mw@lkingw1th477@ng3ls.M.McDonald@faa.gov", "Imwlkingw1th477ng3ls.M.McDonald@faa.gov"},
	}

	for i := 0; i < len(m); i++ {
		s := ScrubEmailAddr(m[i].bad)
		if s != m[i].good {
			t.Errorf("ScrubEmailAddr( %s )  Expect %s, got %s\n", m[i].bad, m[i].good, s)
		}
	}
}
