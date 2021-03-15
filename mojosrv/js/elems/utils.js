"use strict";
"esversion: 8";
/*global
    w2ui, $, app, console, w2utils, document,
    getStatusString, showForm, initializePersonRecord,

*/
// setInnerHTML sets the inner html for the supplied id.  If the id is not
//              found then no action is taken.
//
// INPUTS:
// id = id of element for which HTML will be set
// s  = string containing the HTML
//
// RETURNS:
// 0 = success
// 1 = did not find the label
//------------------------------------------------------------------------------
function setInnerHTML(id,s) {
    var x = document.getElementById(id);
    if (x == null) {
        return 1;
    }
    x.innerHTML = s;
}
