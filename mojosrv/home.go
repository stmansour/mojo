package main

import (
	"fmt"
	"html/template"
	"net/http"
	"mojo/util"
)

// MojoUISupport is a structure of data that will be passed to all html pages.
// It is the responsibility of the page function to populate the data needed by
// the page. The recommendation is to populate only the data needed.
type MojoUISupport struct {
	Language string // what language
	Template string // which template
	ErrMsg   string
}

// HomeUIHandler sends the main UI to the browser
// The forms of the url that are acceptable:
//		/home/
//		/home/<lang>
//		/home/<lang>/<tmpl>
//
// <lang> specifies the language.  The default is en-us
// <tmpl> specifies which template to use. The default is "dflt"
//------------------------------------------------------------------
func HomeUIHandler(w http.ResponseWriter, r *http.Request) {
	var ui MojoUISupport
	var err error
	funcname := "HomeUIHandler"
	appPage := "home.html"
	lang := "en-us"
	tmpl := "default"

	ui.Language = lang
	ui.Template = tmpl

	t, err := template.New(appPage).ParseFiles("./html/" + appPage)
	if nil != err {
		s := fmt.Sprintf("%s: error loading template: %v\n", funcname, err)
		ui.ErrMsg += s
		fmt.Println(s)
	}

	err = t.Execute(w, &ui)

	if nil != err {
		util.LogAndPrintError(funcname, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
