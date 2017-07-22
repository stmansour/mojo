package main

import (
	"mojo/util"
	"net/http"
)

// SvcDisableConsole disables console messages from printing out
func SvcDisableConsole(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.DisableConsole()
	SvcWriteSuccessResponse(w)
}

// SvcEnableConsole enables console messages to print out
func SvcEnableConsole(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.EnableConsole()
	SvcWriteSuccessResponse(w)
}
