package util

import (
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// NonASCIIToASCII maps non-ascii characters to near-ascii-equivalent characters
//
// INPUTS
//  s = string to be converted
//
// RETURNS
//  string = ascii equivalent string
//---------------------------------------------------------------------------------
func NonASCIIToASCII(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}
