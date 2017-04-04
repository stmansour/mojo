package main

import (
	"fmt"
	"net/http"
)

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
	fmt.Printf("Need to implement HomeUIHandler\n")
	fmt.Fprintf(w, "Will implement HomeUIHandler soon!")
	// var ui RRuiSupport
	// var err error
	// funcname := "HomeUIHandler"
	// appPage := "home.html"
	// lang := "en-us"
	// tmpl := "default"

	// path := "/home/"                // this is the part of the URL that got us into this handler
	// uri := r.RequestURI[len(path):] // this pulls off the specific request

	// s, err := url.QueryUnescape(strings.TrimSpace(r.URL.String()))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Printf("HOME HANDLER:  RL = %s\n", s)

	// f := rlib.Stripchars(r.FormValue("filename"), `"`)
	// if len(f) > 0 {

	// 	appPage = strings.TrimSpace(f)
	// }

	// // use   http://domain/home/{lang}/{tmpl}  to set template
	// if len(uri) > 0 {
	// 	s1 := strings.Split(uri, "?")
	// 	sa := strings.Split(s1[0], "/")
	// 	n := len(sa)
	// 	if n > 0 {
	// 		lang = sa[0]
	// 		if n > 1 {
	// 			tmpl = sa[1]
	// 		}
	// 	}
	// }

	// ui.Language = lang
	// ui.Template = tmpl

	// t, err := template.New(appPage).Funcs(RRfuncMap).ParseFiles("./html/" + appPage)
	// if nil != err {
	// 	s := fmt.Sprintf("%s: error loading template: %v\n", funcname, err)
	// 	ui.ReportContent += s
	// 	fmt.Print(s)
	// }
	// err = t.Execute(w, &ui)

	// if nil != err {
	// 	rlib.LogAndPrintError(funcname, err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}
